package torrent

import (
	"io"
	"net"

	"github.com/anacrolix/dht/v2"
	"github.com/anacrolix/dht/v2/krpc"
)

type DhtServer interface {
	Stats() interface{}
	ID() [20]byte
	Addr() net.Addr
	AddNode(ni krpc.NodeInfo) error
	Ping(addr *net.UDPAddr)
	Announce(hash [20]byte, port int, impliedPort bool) (DhtAnnounce, error)
	WriteStatus(io.Writer)
}

type BEP44 interface {
	AddStorageItem(item StorageItem) bool
	GetStorageItem(target [20]byte) (StorageItem, bool)
	ArbitraryData(target [20]byte, seq *uint64) (ArbitraryData, error)
}

type DhtAnnounce interface {
	Close()
	Peers() <-chan dht.PeersValues
}

type StorageItem interface{}

type ArbitraryData interface{}

type anacrolixDhtServerWrapper struct {
	*dht.Server
}

func (me anacrolixDhtServerWrapper) Stats() interface{} {
	return me.Server.Stats()
}

type anacrolixDhtAnnounceWrapper struct {
	*dht.Announce
}

func (me anacrolixDhtAnnounceWrapper) Peers() <-chan dht.PeersValues {
	return me.Announce.Peers
}

func (me anacrolixDhtServerWrapper) Announce(hash [20]byte, port int, impliedPort bool) (DhtAnnounce, error) {
	ann, err := me.Server.Announce(hash, port, impliedPort)
	return anacrolixDhtAnnounceWrapper{ann}, err
}

func (me anacrolixDhtServerWrapper) Ping(addr *net.UDPAddr) {
	me.Server.Ping(addr, nil)
}

func (me anacrolixDhtServerWrapper) AddStorageItem(i StorageItem) bool {
	return me.Server.AddStorageItem(i.(dht.StorageItem))
}

func (me anacrolixDhtServerWrapper) GetStorageItem(k [20]byte) (StorageItem, bool) {
	return me.Server.GetStorageItem(k)
}

func (me anacrolixDhtServerWrapper) ArbitraryData(k [20]byte, seq *uint64) (ArbitraryData, error) {
	return me.Server.ArbitraryData(k, seq)
}

var _ DhtServer = anacrolixDhtServerWrapper{}
var _ BEP44 = anacrolixDhtServerWrapper{}
