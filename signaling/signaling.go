package signalling

import (
	"fmt"
	"sync"
	"time"

	grpc "github.com/pigeatgarlic/signaling/gRPC"
	"github.com/pigeatgarlic/signaling/protocol"
	"github.com/pigeatgarlic/signaling/websocket"
	"github.com/pigeatgarlic/webrtc-proxy/signalling/gRPC/packet"
)


type Signalling struct {
	waitLine map[string]*WaitingTenant
	pairs    map[string]*Pair
	mut 	 *sync.Mutex
	
	handlers []protocol.ProtocolHandler
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

func (signaling *Signalling)tokenMatch(token string, tent protocol.Tenant) (client protocol.Tenant, worker protocol.Tenant, err error){
	signaling.mut.Lock()
	for index,wait := range signaling.waitLine{
		if index == "client" && token == "server" {
			fmt.Printf("match\n");
			client = wait.waiter
			worker = tent
			return;
		} else if token == "server" && index == "client" {
			fmt.Printf("match\n");
			worker = wait.waiter
			client = tent
			return;
		} else {
			continue
		}
	}
	signaling.mut.Unlock()
	err = fmt.Errorf("no match peer")
	return
}

func InitSignallingServer(conf *protocol.SignalingConfig) *Signalling {
	var err error
	var signaling Signalling
	signaling.handlers = make([]protocol.ProtocolHandler, 2)
	signaling.pairs    = make(map[string]*Pair)
	signaling.waitLine = make(map[string]*WaitingTenant)
	signaling.mut = &sync.Mutex{}

	signaling.handlers[0] = grpc.InitSignallingServer(conf);
	signaling.handlers[1] = ws.InitSignallingWs(conf);
	if err != nil  {
		fmt.Printf("%s\n",err.Error())
		return nil;
	}

	fun := func (token string, tent protocol.Tenant) error {
		client, worker, err := signaling.tokenMatch(token,tent);

		if err == nil{
			pair := &Pair{
				client: client,
				worker: worker,
			}
			signaling.addPair(token,pair)
			pair.handlePair()

			pair.worker.Send(&packet.UserResponse{
				Id: 0,	
				Error: "",
				Data: map[string]string{
					"Target": "START",
				},
			})
		} else {
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
				signaling.removePair(i)
			}
			time.Sleep(10*time.Millisecond);
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


