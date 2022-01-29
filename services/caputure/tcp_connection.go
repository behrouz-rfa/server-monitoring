package capture

import (
	"time"
)

const RECV_TIMEOUT = 2000 * time.Millisecond

type connectionState int

const (
	stateReciving connectionState = iota
	stateClosed
)

type tcpConnection struct {
	ID      string
	ReqID   string
	srcAddr string
	srcPort int
	dstAddr string
	dstPort int
	State   connectionState

	timer            *time.Timer // Used for expire check
	chRecvTimeoutMsg chan *tcpConnection
}

func NewTcpConnection(id string, reqId string, srcAddr string, srcPort int, dstAddr string, dstPort int) *tcpConnection {
	t := &tcpConnection{
		ID:      id,
		ReqID:   reqId,
		State:   stateReciving,
		srcAddr: srcAddr,
		srcPort: srcPort,
		dstAddr: dstAddr,
		dstPort: dstPort,
	}
	t.timer = time.AfterFunc(RECV_TIMEOUT, t.recvTimeout)
	return t
}

func (t *tcpConnection) RegisterTimeoutEvent(chRecvTimeoutMsg chan *tcpConnection) {
	t.chRecvTimeoutMsg = chRecvTimeoutMsg
}

func (t *tcpConnection) IsRequest(srcAddr string, srtPort int) bool {
	return t.srcAddr == srcAddr && t.srcPort == srtPort
}

func (t *tcpConnection) IsSameRequest(reqId string) bool {
	return t.ReqID == reqId
}

func (t *tcpConnection) Close() {
	t.State = stateClosed
}

func (t *tcpConnection) MarkActive() {
	t.timer.Reset(RECV_TIMEOUT)
}

func (t *tcpConnection) recvTimeout() {
	if t.chRecvTimeoutMsg != nil {
		t.chRecvTimeoutMsg <- t // Notify connection not recieve stream
	}
}
