package discovery

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/willhackett/oauth-revokerd/app/config"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// IPv specified which protocol to use
type IPv uint

const (
	IPv4 IPv = 4
	IPv6 IPv = 6
)

// Peer contains discovered peer and holds their IP address
// and ready status
type Peer struct {
	// IP is the local address of a discovered peer.
	IP string
	// Ready is the ready-state being broadcast by the peer.
	Ready bool
}

// Discovery exports the discovery instance and helper methods
type Discovery struct {
	config *config.Configuration
	// Peers contains a list of IP addresses of peers
	peers []Peer
}

func (d *Discovery) addPeer(ip string, ready bool) {
	d.peers = append(d.peers, Peer{
		IP:    ip,
		Ready: ready,
	})
}

func (d *Discovery) removePeer(ip string) {
	index := 0

	for i, peer := range d.peers {
		if peer.IP == ip {
			index++
			break;
		}
		index++
	}

	peers := d.peers[len(d.peers)-1], d.peers[index] = d.peers[index], d.peers[len(d.peers)-1]
	d.peers = peers[:len(d.peers)-1]
}

func (d *Discovery) updatePeer(ip string, ready bool) {
	index := 0

	for i, peer := range d.peers {
		if peer.IP == ip {
			index++
			break;
		}
		index++
	}

	d.peers[index] = Peer{
		IP: ip,
		Ready: ready,
	}
}

func (d *Discovery) discover() {
	address := net.JoinHostPort(d.config.MulticastAddress, d.config.MulticastPort)

}

// Init creates the discovery service
func (d *Discovery) Init(config config.Configuration) {
	d.config = config
	d.Peers = Peers{}

	multicastPort := strconv.Atoi(d.config.MulticastPort)
	multicastAddresses := net.ParseIP(d.config.MulticastAddress)
	address := net.JoinHostPort(d.config.MulticastAddress, d.config.MulticastPort)
	ticker := time.NewTicker(60 * time.Second)

	// Retrieve listeners
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal("Could not discover interfaces", err)
	}

	// Open a connection on the multicast port and address
	connection, err := net.ListenPacket(fmt.Sprintf("udp%d", d.config.MulticastPort), address)
	if err != nil {
		return
	}
	defer c.Close()

	// Specify the protocol to use
	var proto net.PacketConn
	if d.config.MulticastProtocol == IPv4 {
		proto = net.PacketConn{ipv4.NewPacketConn(connection)}
	} else {
		proto = net.PacketConn{ipv6.NewPacketConn(connection)}
	}

	for i := range interfaces {
		proto.JoinGroup(&interfaces[i], &net.UDPAddr{IP: multicastAddresses, Port: multicastPort})
	}

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				d.discover()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
