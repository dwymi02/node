package p2p

import (
	"encoding/binary"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"net"
	"sync"
	"time"
)

type Peer struct {
	IsPermitCompleteNode bool // block msg

	Name string
	Id   []byte // len 32

	TcpConn net.Conn

	publicIPv4    []byte // public IP if is public peer
	tcpListenPort int
	udpListenPort int

	knownPeerIds                       mapset.Set // set[string([32]byte)] // string key
	knownPeerKnowledgeDuplicateRemoval sync.Map   // map[string]set[string(byte)]

	connTime   time.Time
	activeTime time.Time // check live

}

func NewPeer(id []byte, name string) *Peer {
	return &Peer{
		IsPermitCompleteNode: false,
		Id:                   id,
		Name:                 name,
		publicIPv4:           nil,
		tcpListenPort:        0,
		udpListenPort:        0,
		knownPeerIds:         mapset.NewSet(),
	}
}

func (p *Peer) Close() {
	if p.TcpConn != nil {
		p.TcpConn.Close()
		p.TcpConn = nil
	}
	// reset
	p.IsPermitCompleteNode = false
}

func (p *Peer) ParseRemotePublicTCPAddress() []byte {
	if p.publicIPv4 == nil {
		return nil
	}
	bts := make([]byte, 6)
	copy(bts[0:4], p.publicIPv4)
	binary.BigEndian.PutUint16(bts[4:6], uint16(p.tcpListenPort))
	//fmt.Println("ParseRemotePublicTCPAddress", bts)
	return bts
}

func (p *Peer) AddKnownPeerId(pid []byte) {
	p.knownPeerIds.Add(string(pid))
	if p.knownPeerIds.Cardinality() > 200 {
		p.knownPeerIds.Pop() // remove one
	}
}

func (p *Peer) SendMsg(ty uint16, msgbody []byte) error {
	if msgbody == nil {
		msgbody = []byte{}
	}
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, ty)
	var dtlen []byte = nil
	if ty == TCPMsgTypeData {
		dtlen = make([]byte, 4)
		binary.BigEndian.PutUint32(dtlen, uint32(len(msgbody)))
	} else {
		if len(msgbody) > 65535 {
			return fmt.Errorf("msg body size overflow 65535.")
		}
		dtlen = make([]byte, 2)
		binary.BigEndian.PutUint16(dtlen, uint16(len(msgbody)))
	}
	data = append(data, dtlen...)
	data = append(data, msgbody...)
	// send data
	if p.TcpConn != nil {
		_, e := p.TcpConn.Write(data)
		if e != nil {
			fmt.Println("SendMsg error", e)
			return e
		}
	}
	return nil

}
