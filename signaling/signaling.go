package signalling

import (
	"fmt"
	"sync"
	"time"

	grpc "github.com/pigeatgarlic/signaling/gRPC"
	"github.com/pigeatgarlic/signaling/protocol"
	"github.com/pigeatgarlic/signaling/validator"
	"github.com/pigeatgarlic/signaling/websocket"
	"github.com/pigeatgarlic/webrtc-proxy/signalling/gRPC/packet"
)


type Signalling struct {
	waitLine map[validator.ValidationResult]*WaitingTenant
	pairs    map[int]*Pair
	mut 	 *sync.Mutex
	
	handlers []protocol.ProtocolHandler
	validator validator.Validator
}

func (signaling *Signalling)removePair(s int) {
	signaling.mut.Lock()
	delete(signaling.pairs,s)
	signaling.mut.Unlock()
}
func (signaling *Signalling)removeTenant(s validator.ValidationResult) {
	signaling.mut.Lock()
	delete(signaling.waitLine,s)
	signaling.mut.Unlock()
}
func (signaling *Signalling)addPair(s int, tenant *Pair) {
	signaling.mut.Lock()
	signaling.pairs[s] = tenant;
	signaling.mut.Unlock()
}
func (signaling *Signalling)addTenant(s validator.ValidationResult, tenant *WaitingTenant) {
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
	if req.Target == "ICE" || req.Target == "SDP" || req.Target == "START" || req.Target == "PREFLIGHT"{
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
func (signaling *Signalling)tokenMatch(result validator.ValidationResult, tent protocol.Tenant) (client protocol.Tenant, worker protocol.Tenant, found bool,id int){
	found = true;
	signaling.mut.Lock()
	defer func ()  { signaling.mut.Unlock() }()


	for index,wait := range signaling.waitLine{
		if index.ID == result.ID && (result.IsServer == !index.IsServer){
			fmt.Printf("match\n");
			if result.IsServer {
				client = wait.waiter
				worker = tent
			} else {
				worker = wait.waiter
				client = tent
			}
			return;
		} else {
			continue;
		}
	}
	found = false
	return
}

func InitSignallingServer(conf *protocol.SignalingConfig, provider validator.Validator) *Signalling {
	var signaling Signalling
	signaling.handlers = make([]protocol.ProtocolHandler, 2)
	signaling.pairs    = make(map[int]*Pair)
	signaling.waitLine = make(map[validator.ValidationResult]*WaitingTenant)
	signaling.mut = &sync.Mutex{}

	signaling.handlers[0] = grpc.InitSignallingServer(conf);
	signaling.handlers[1] = ws.InitSignallingWs(conf);
	signaling.validator = provider

	fun := func (token string, tent protocol.Tenant) error {
		result,err := signaling.validator.Validate(token);
		if err != nil {
			fmt.Printf("validation error %s\n",err.Error())
			return err;
		}
		client, worker, found, id := signaling.tokenMatch(*result,tent);

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
			signaling.addTenant(*result,wait)
			wait.handle()
		}
		return nil;
	};

	go func() {
		for {
			var rev []int;
			
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
			var rev []validator.ValidationResult;

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


