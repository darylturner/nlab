package network

import (
	"errors"
	"fmt"
	"github.com/darylturner/nlab/internal/config"
	"github.com/darylturner/nlab/internal/network/netlinux"
	"math/rand"
	"runtime"
	"time"
)

type Network interface {
	Create() error
	Destroy() error
	AddNode(string)
}

type PseudoWire struct {
	Type     string
	Init     bool
	BasePort int
}

func GetMap(cfg *config.Topology) (map[string]Network, error) {
	allLinks := make(map[string]Network)
	os := runtime.GOOS
	for _, nd := range cfg.Nodes {
		if nd.Network.Management == true {
			link := "_" + cfg.ManagementBridge
			nd.Network.Links = append(nd.Network.Links, link)
		}

		for _, link := range nd.Network.Links {
			if net, ok := allLinks[link]; ok {
				net.AddNode(nd.Tag)
				continue
			}

			switch os {
			case "linux":
				net := netlinux.New(link)
				net.AddNode(nd.Tag)
				allLinks[link] = &net
			default:
				return allLinks, errors.New("running on unsupported os")
			}
		}
	}

	return allLinks, nil
}

func GetPseudoWireMap(cfg *config.Topology) (map[string]*PseudoWire, error) {
	allLinks := make(map[string]*PseudoWire)
	portBase := 30000
	os := runtime.GOOS
	for _, nd := range cfg.Nodes {
		for _, link := range nd.Network.Links {
			if _, ok := allLinks[link]; ok {
				continue
			}

			switch os {
			case "linux":
				allLinks[link] = &PseudoWire{
					Type:     "qemu-unicast-udp",
					Init:     false,
					BasePort: portBase,
				}
			default:
				return allLinks, errors.New("running on unsupported os")
			}

			portBase += 2
		}
	}

	return allLinks, nil
}

func randomByte() int {
	return rand.Intn(256)
}

func RandomMAC() string {
	baseMAC := "52:54:00"
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(
		baseMAC+":%02x:%02x:%02x",
		randomByte(),
		randomByte(),
		randomByte(),
	)
}
