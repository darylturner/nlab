package netlinux

import (
	"crypto/sha256"
	"fmt"
	"os/exec"
	"strings"
)

func TapUID(link, node string) string {
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
	tag := brTap.Tag

	// if bridge is prefixed with _ skip the create and strip
	// the prefix
	if !strings.HasPrefix(tag, "_") {
		if err := createBridge(tag); err != nil {
			return err
		}
	} else {
		tag = strings.TrimLeft(tag, "_")
	}

	for _, node := range brTap.NodesJoined {
		tapName := TapUID(tag, node)
		if err := createTap(tapName, tag); err != nil {
			return err
		}
	}

	return nil
}

func (brTap BridgeTapNetwork) Destroy() error {
	tag := brTap.Tag

	// if bridge is prefixed with _ skip the destroy and strip
	// the prefix
	if !strings.HasPrefix(tag, "_") {
		if err := destroyBridge(tag); err != nil {
			return err
		}
	} else {
		tag = strings.TrimLeft(tag, "_")
	}

	for _, node := range brTap.NodesJoined {
		tapName := TapUID(tag, node)
		if err := destroyTap(tapName); err != nil {
			return err
		}
	}

	return nil
}

func createTap(name, bridge string) error {
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
	if err := exec.Command("ip", "link", "add", name, "type", "bridge").Run(); err != nil {
		return err
	}
	if err := exec.Command("ip", "link", "set", name, "up").Run(); err != nil {
		return err
	}

	return nil
}

func destroyTap(name string) error {
	if err := exec.Command("ip", "link", "delete", name).Run(); err != nil {
		return err
	}

	return nil
}

func destroyBridge(name string) error {
	if err := exec.Command("ip", "link", "delete", name, "type", "bridge").Run(); err != nil {
		return err
	}

	return nil
}
