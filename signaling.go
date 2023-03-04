package signaling

import (
	"sync"
	"time"

	"github.com/pigeatgarlic/signaling/protocol/websocket"
	grpc "github.com/pigeatgarlic/signaling/protocol/gRPC"

	"github.com/pigeatgarlic/signaling/protocol"
	"github.com/pigeatgarlic/signaling/validator"
)

type Pair struct {
	client protocol.Tenant
	worker protocol.Tenant
}

type Signalling struct {
	waitLine map[string]protocol.Tenant
	pairs    map[int64]Pair
	mut      *sync.Mutex

	handlers  []protocol.ProtocolHandler
	validator validator.Validator
}

func InitSignallingServer(conf *protocol.SignalingConfig, provider validator.Validator) *Signalling {
	signaling := Signalling {
		waitLine : make(map[string]protocol.Tenant),
		pairs : make(map[int64]Pair),
		mut : &sync.Mutex{},
		validator : provider,
		handlers : []protocol.ProtocolHandler{
			grpc.InitSignallingServer(conf),
			ws.InitSignallingWs(conf),
		},
	}


	go func() {
		for {
			var rev []int64
			signaling.mut.Lock()
			for index, pair := range signaling.pairs {
				if pair.client.IsExited() {
					pair.worker.Exit()
					rev = append(rev, index)
				} else if pair.worker.IsExited() {
					pair.client.Exit()
					rev = append(rev, index)
				}
			}
			signaling.mut.Unlock()

			for _, i := range rev {
				signaling.removePair(i)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	go func() {
		for {
			var rev []string
			signaling.mut.Lock()
			for index, wait := range signaling.waitLine {
				if wait.IsExited() {
					rev = append(rev, index)
				}
			}
			signaling.mut.Unlock()
			for _, i := range rev {
				signaling.removeTenant(i)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	go func ()  {
		for {
			for _,t := range signaling.waitLine {
				if t.Peek() {
					_ = t.Receive() // discard
				}
			}
		}
	}()

	for _, handler := range signaling.handlers {
		handler.OnTenant(func(token string, tent protocol.Tenant) error {
			signaling.addTenant(token,tent)

			keys := make([]string, 0, len(signaling.waitLine))
			for k := range signaling.waitLine {
				keys = append(keys, k)
			}

			pairs,new_queue := signaling.validator.Validate(keys)

			for _,k := range keys {
				rm := true
				for _,n := range new_queue {
					if n == k {
						rm = false
					}
				}

				if rm {
					delete(signaling.waitLine,k)
				}
			}


			for k,v := range pairs {
				pair := Pair {client: nil,worker: nil}
				for _, v2 := range keys {
					if v2 == k {
						pair.worker = signaling.waitLine[v2]
					} else if v2 == v {
						pair.client = signaling.waitLine[v2]
					}
				}

				
				if pair.client == nil || pair.worker == nil {
					continue
				}

				signaling.addPair(time.Now().UnixMicro(),pair)
			}

			return nil
		})
	}

	return &signaling
}



func (pair *Pair) handlePair() {
	go func() {
		for {
			if pair.client.IsExited() || pair.worker.IsExited() { return }
			pair.client.Send(pair.worker.Receive())
		}
	}()
	go func() {
		for {
			if pair.client.IsExited() || pair.worker.IsExited() { return }
			pair.worker.Send(pair.client.Receive())
		}
	}()
}
