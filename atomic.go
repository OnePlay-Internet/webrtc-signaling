package signaling

import "github.com/pigeatgarlic/signaling/protocol"

func (signaling *Signalling) removePair(s int64) {
	signaling.mut.Lock()
	delete(signaling.pairs, s)
	signaling.mut.Unlock()
}
func (signaling *Signalling) addPair(s int64, tenant Pair) {
	signaling.mut.Lock()
	signaling.pairs[s] = tenant
	signaling.mut.Unlock()
}

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