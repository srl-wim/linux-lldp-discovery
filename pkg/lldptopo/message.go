package lldptopo

type msgID int

const (
	// Used to request to shutdown, expects ACK on respChan
	shutdown msgID = iota
	// Acknowledge a request.
	ack
	// triggers link status validation
	doWork
	
)

// Control message channel type
type cMsg struct {
	id       msgID
	data     []byte
	respChan chan *cMsg
}