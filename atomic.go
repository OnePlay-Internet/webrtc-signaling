package signaling

import "github.com/thinkonmay/signaling-server/protocol"

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
