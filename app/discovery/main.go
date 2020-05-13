package discovery

import (
	"context"
	"fmt"

	"github.com/brutella/dnssd"
	"github.com/willhackett/oauth-revokerd/app/config"
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

// func (d *Discovery) removePeer(ip string) {
// 	index := 0

// 	for i, peer := range d.peers {
// 		if peer.IP == ip {
// 			index++
// 			break
// 		}
// 		index++
// 	}

// 	peers := d.peers[len(d.peers)-1], d.peers[index] == d.peers[index], d.peers[len(d.peers)-1]
// 	d.peers = peers[:len(d.peers)-1]
// }

func (d *Discovery) updatePeer(ip string, ready bool) {
	index := 0

	for _, peer := range d.peers {
		if peer.IP == ip {
			index++
			break
		}
		index++
	}

	d.peers[index] = Peer{
		IP:    ip,
		Ready: ready,
	}
}

// func (d *Discovery) discover() {
// 	address := net.JoinHostPort(d.config.MulticastAddress, d.config.MulticastPort)

// }

// Init starts announcing the service is ready on the network
func (d *Discovery) Init(config *config.Configuration) {
	d.config = config
	d.peers = []Peer{}

	cfg := dnssd.Config{
		Name: "oauth-revokerd",
		Type: "_http._tcp",
		Port: config.Port,
	}
	svc, _ := dnssd.NewService(cfg)
	resp, _ := dnssd.NewResponder()

	hdl, _ := resp.Add(svc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp.Respond(ctx)

	hdl.UpdateText(map[string]string{"ready": "true"}, resp)

	fmt.Println(svc)
}
