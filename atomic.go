package signaling

import "github.com/pigeatgarlic/signaling/protocol"


func (signaling *Signalling) removeTenant(s string) {
	signaling.mut.Lock()
	delete(signaling.waitLine, s)
	signaling.mut.Unlock()
}

func (signaling *Signalling) addTenant(s string, tenant protocol.Tenant) {
	signaling.mut.Lock()
	signaling.waitLine[s] = tenant
	signaling.mut.Unlock()
}