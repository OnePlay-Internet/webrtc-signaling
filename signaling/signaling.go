package signalling

import (
	"fmt"
	"sync"
	"time"

	grpc "github.com/pigeatgarlic/signaling/gRPC"
	"github.com/pigeatgarlic/signaling/protocol"
	"github.com/pigeatgarlic/signaling/validator"
	"github.com/pigeatgarlic/signaling/validator/oneplay"
	"github.com/pigeatgarlic/signaling/websocket"
	"github.com/pigeatgarlic/webrtc-proxy/signalling/gRPC/packet"
)


type Signalling struct {
	waitLine map[string]*WaitingTenant
	pairs    map[string]*Pair
	mut 	 *sync.Mutex
	
	handlers []protocol.ProtocolHandler
	validator validator.Validator
}

func (signaling *Signalling)removePair(s string) {
	signaling.mut.Lock()
	delete(signaling.pairs,s)
	signaling.mut.Unlock()
}
func (signaling *Signalling)removeTenant(s string) {
	signaling.mut.Lock()
	delete(signaling.waitLine,s)
	signaling.mut.Unlock()
}
func (signaling *Signalling)addPair(s string, tenant *Pair) {
	signaling.mut.Lock()
	signaling.pairs[s] = tenant;
	signaling.mut.Unlock()
}
func (signaling *Signalling)addTenant(s string, tenant *WaitingTenant) {
	signaling.mut.Lock()
	signaling.waitLine[s] = tenant;
	signaling.mut.Unlock()
}

func ProcessReq(req *packet.UserRequest)*packet.UserResponse  {
	if req == nil {
		return nil;			
	}

	var res packet.UserResponse; 
	res.Data = req.Data
	if req.Target == "ICE" || req.Target == "SDP" || req.Target == "START"{
		fmt.Printf("forwarding %s packet\n",req.Target);
		res.Data["Target"] = req.Target;
	} else {
		fmt.Printf("unknown %s packet\n",req.Target);
		return nil;
	}

	return &res;
}



type WaitingTenant struct {
	waiter protocol.Tenant
}

func (wait *WaitingTenant) handle(){
	go func() {
		for { 
			if wait.waiter.IsExited() {
				return;
			}
			time.Sleep(time.Millisecond)
		}	
	}()
}

type Pair struct {
	client protocol.Tenant
	worker protocol.Tenant
}

func (pair *Pair) handlePair(){
	go func ()  {
		for {
			if pair.client.IsExited() || pair.worker.IsExited(){
				return;	
			}
			dat := pair.worker.Receive()	
			pair.client.Send(ProcessReq(dat))
			time.Sleep(time.Millisecond)
		}	
	}()

	go func ()  {
		for {
			if pair.client.IsExited() || pair.worker.IsExited(){
				return;	
			}
			dat := pair.client.Receive()	
			pair.worker.Send(ProcessReq(dat))
			time.Sleep(time.Millisecond)
		}	
	}()
}

// TODO
func (signaling *Signalling)tokenMatch(token string, tent protocol.Tenant) (client protocol.Tenant, worker protocol.Tenant, found bool, id string){
	found = true;
	signaling.mut.Lock()
	defer func ()  { signaling.mut.Unlock() }()

	result := signaling.validator.Validate(token);

	for index,wait := range signaling.waitLine{
		if index == result.ClientToken && token == result.ServerToken {
			fmt.Printf("match\n");
			client = wait.waiter
			worker = tent
			return;
		} else if token == result.ClientToken && index == result.ServerToken {
			fmt.Printf("match\n");
			worker = wait.waiter
			client = tent
			return;
		} else if token == index {
			wait.waiter.Exit()
		} else {
			continue;
		}
	}
	found = false
	return
}

func InitSignallingServer(conf *protocol.SignalingConfig) *Signalling {
	var signaling Signalling
	signaling.handlers = make([]protocol.ProtocolHandler, 2)
	signaling.pairs    = make(map[string]*Pair)
	signaling.waitLine = make(map[string]*WaitingTenant)
	signaling.mut = &sync.Mutex{}

	signaling.handlers[0] = grpc.InitSignallingServer(conf);
	signaling.handlers[1] = ws.InitSignallingWs(conf);
	signaling.validator = oneplay.NewOneplayValidator(conf.ValidationUrl)

	fun := func (token string, tent protocol.Tenant) error {
		client, worker, found, id := signaling.tokenMatch(token,tent);

		if found {
			fmt.Printf("new pair\n")
			pair := &Pair{
				client: client,
				worker: worker,
			}
			signaling.addPair(id,pair)
			pair.handlePair()

			pair.worker.Send(&packet.UserResponse{
				Id: 0,	
				Error: "",
				Data: map[string]string{
					"Target": "START",
				},
			})
		} else {
			fmt.Printf("new tenant to waitline\n")
			wait := &WaitingTenant{
				waiter: tent,
			};
			signaling.addTenant(token,wait)
			wait.handle()
		}
		return nil;
	};

	go func() {
		for {
			var rev []string;
			
			signaling.mut.Lock()
			for index, pair := range signaling.pairs{
				if  pair.client.IsExited(){
					pair.worker.Exit()
					rev = append(rev, index);	
				} else if pair.worker.IsExited(){
					pair.client.Exit()
					rev = append(rev, index);	
				}
			}
			signaling.mut.Unlock()

			for _,i := range rev {
				fmt.Printf("removing pair\n");
				signaling.removePair(i)
			}
			time.Sleep(10*time.Millisecond);
		}	
	}()

	go func() {
		for {
			signaling.mut.Lock()
			for i,_ := range signaling.pairs {
				fmt.Printf("pair, holding id %s\n",i);	
			}
			for i,_ := range signaling.waitLine {
				fmt.Printf("tenant in waiting line, holding token %s\n",i);	
			}
			signaling.mut.Unlock()
			time.Sleep(10*time.Second);
		}	
	}()

	go func() {
		for {
			var rev []string;

			signaling.mut.Lock()
			for index, wait := range signaling.waitLine{
				if wait.waiter.IsExited() {
					rev = append(rev, index);	
				}
			}
			signaling.mut.Unlock()

			for _,i := range rev {
				fmt.Printf("removing tenant from waiting line\n");
				signaling.removeTenant(i)
			}
			time.Sleep(10*time.Millisecond);
		}	
	}()

	for _,handler := range signaling.handlers {
		handler.OnTenant(fun);
	}
	return &signaling;
}


