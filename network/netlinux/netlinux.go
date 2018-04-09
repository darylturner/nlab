package netlinux

import (
	"crypto/sha256"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func tapUID(link, node string) string {
	// required as linux has a character limit on tap names,
	// we need to be fairly sure of deriving a unique tap name
	// to avoid name collisions.
	hash := sha256.New()
	hash.Write([]byte(link))
	hash.Write([]byte(node))

	return fmt.Sprintf("veth%.5x", hash.Sum(nil))
}

func New(tag string) BridgeTapNetwork {
	return BridgeTapNetwork{Tag: tag}
}

type BridgeTapNetwork struct {
	Tag         string   `json:"segment"`
	NodesJoined []string `json:"nodes"`
}

func (brTap *BridgeTapNetwork) AddNode(node string) {
	brTap.NodesJoined = append(brTap.NodesJoined, node)

	return
}

func (brTap BridgeTapNetwork) Create() error {
	if err := createBridge(brTap.Tag); err != nil {
		return err
	}

	for _, node := range brTap.NodesJoined {
		tapName := tapUID(brTap.Tag, node)
		if err := createTap(tapName, brTap.Tag); err != nil {
			return err
		}
	}

	return nil
}

func (brTap BridgeTapNetwork) Destroy() error {
	if err := destroyBridge(brTap.Tag); err != nil {
		return err
	}

	for _, node := range brTap.NodesJoined {
		tapName := tapUID(brTap.Tag, node)
		if err := destroyTap(tapName); err != nil {
			return err
		}
	}

	return nil
}

func createTap(name, bridge string) error {
	contextLog := log.WithFields(log.Fields{
		"tap": name, "bridge": bridge,
	})

	contextLog.Info("creating tap")
	if err := exec.Command("ip", "tuntap", "add", "dev", name, "mode", "tap").Run(); err != nil {
		return err
	}
	if err := exec.Command("ip", "link", "set", name, "master", bridge).Run(); err != nil {
		return err
	}
	if err := exec.Command("ip", "link", "set", name, "up").Run(); err != nil {
		return err
	}

	return nil
}

func createBridge(name string) error {
	contextLog := log.WithFields(log.Fields{
		"bridge": name,
	})

	contextLog.Info("creating bridge")
	if err := exec.Command("ip", "link", "add", name, "type", "bridge").Run(); err != nil {
		return err
	}
	if err := exec.Command("ip", "link", "set", name, "up").Run(); err != nil {
		return err
	}

	return nil
}

func destroyTap(name string) error {
	contextLog := log.WithFields(log.Fields{
		"tap": name,
	})

	contextLog.Info("destroying tap")
	if err := exec.Command("ip", "link", "delete", name).Run(); err != nil {
		return err
	}

	return nil
}

func destroyBridge(name string) error {
	contextLog := log.WithFields(log.Fields{
		"bridge": name,
	})

	contextLog.Info("destroying bridge")
	if err := exec.Command("ip", "link", "delete", name, "type", "bridge").Run(); err != nil {
		return err
	}

	return nil
}
