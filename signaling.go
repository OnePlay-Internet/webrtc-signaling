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


type Pair struct {
	client protocol.Tenant
	worker protocol.Tenant
}

type Signalling struct {
	waitLine map[string]protocol.Tenant
	pairs    map[int]Pair
	mut 	 *sync.Mutex
	
	handlers []protocol.ProtocolHandler
	validator validator.Validator
}

func (signaling *Signalling)removePair(s int) {
	signaling.mut.Lock()
	delete(signaling.pairs,s)
	signaling.mut.Unlock()
}
func (signaling *Signalling)addPair(s int, tenant Pair) {
	signaling.mut.Lock()
	signaling.pairs[s] = tenant;
	signaling.mut.Unlock()
}


func (signaling *Signalling)removeTenant(s string) {
	signaling.mut.Lock()
	delete(signaling.waitLine,s)
	signaling.mut.Unlock()
}

func (signaling *Signalling)addTenant(s string, tenant protocol.Tenant) {
	signaling.mut.Lock()
	signaling.waitLine[s] = tenant;
	signaling.mut.Unlock()
}
func InitSignallingServer(conf *protocol.SignalingConfig, provider validator.Validator) *Signalling {
	var signaling Signalling
	signaling.waitLine = make(map[string]protocol.Tenant)
	signaling.pairs    = make(map[int]Pair)
	signaling.mut = &sync.Mutex{}

	signaling.handlers = []protocol.ProtocolHandler{
		grpc.InitSignallingServer(conf),
		ws.InitSignallingWs(conf),
	}

	signaling.validator = provider

	fun := func (token string, tent protocol.Tenant) error {
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
				if wait.IsExited() { rev = append(rev, index);	}
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














func ProcessReq(req *packet.UserRequest)*packet.UserResponse  {
	if req == nil {
		return nil;			
	}

	var res packet.UserResponse; 
	res.Data = req.Data
	if req.Target == "ICE" || req.Target == "SDP" || req.Target == "START" || req.Target == "PREFLIGHT" || req.Target == "NVCOMPUTER" || req.Target == "SELECTION" || req.Target == "RESPONSE"{
		fmt.Printf("forwarding %s packet\n",req.Target);
		res.Data["Target"] = req.Target;
	} else {
		fmt.Printf("unknown %s packet\n",req.Target);
		return nil;
	}

	return &res;
}




func (wait *protocol.Tenant) handle(){
	go func() {
		for { 
			if wait.IsExited() {
				return;
			}
			time.Sleep(time.Millisecond)
		}	
	}()
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
