package netlinux

import (
	"crypto/sha256"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func TapUID(id int, link, node string) string {
	// required as linux has a character limit on tap names,
	// we need to be fairly sure of deriving a unique tap name
	// to avoid name collisions.
	hash := sha256.New()
	hash.Write([]byte(strconv.Itoa(id)))
	hash.Write([]byte(link))
	hash.Write([]byte(node))

	return fmt.Sprintf("veth%.5x", hash.Sum(nil))
}

func New(id int, tag string) BridgeTapNetwork {
	return BridgeTapNetwork{LabID: id, Tag: tag}
}

type BridgeTapNetwork struct {
	LabID       int      `json:"lab_id"`
	Tag         string   `json:"segment"`
	NodesJoined []string `json:"nodes"`
}

func (brTap *BridgeTapNetwork) AddNode(node string) {
	brTap.NodesJoined = append(brTap.NodesJoined, node)

	return
}

func (brTap BridgeTapNetwork) Create() error {
	tag := brTap.Tag
	id := brTap.LabID

	// if bridge is prefixed with _ skip the create and strip
	// the prefix
	if !strings.HasPrefix(tag, "_") {
		if err := createBridge(tag); err != nil {
			return err
		}
	} else {
		tag = strings.TrimLeft(tag, "_")
	}

	var err error
	for _, node := range brTap.NodesJoined {
		tapName := TapUID(id, tag, node)
		err = createTap(tapName, tag)
	}
	if err != nil {
		return err
	}

	return nil
}

func (brTap BridgeTapNetwork) Destroy() error {
	tag := brTap.Tag
	id := brTap.LabID

	// if bridge is prefixed with _ skip the destroy and strip
	// the prefix
	if !strings.HasPrefix(tag, "_") {
		if err := destroyBridge(tag); err != nil {
			return err
		}
	} else {
		tag = strings.TrimLeft(tag, "_")
	}

	var err error
	for _, node := range brTap.NodesJoined {
		tapName := TapUID(id, tag, node)
		err = destroyTap(tapName)
	}
	if err != nil {
		return err
	}

	return nil
}

func createTap(name, bridge string) error {
	if err := exec.Command("ip", "tuntap", "add", "dev", name, "mode", "tap").Run(); err != nil {
		return fmt.Errorf("error creating tap %v: %v", name, err)
	}
	if err := exec.Command("ip", "link", "set", name, "master", bridge).Run(); err != nil {
		return fmt.Errorf("error linking tap %v to bridge %v: %v", name, bridge, err)
	}
	if err := exec.Command("ip", "link", "set", name, "up").Run(); err != nil {
		return fmt.Errorf("error setting tap to up %v: %v", name, err)
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
		return fmt.Errorf("error removing tap %v: %v", name, err)
	}

	return nil
}

func destroyBridge(name string) error {
	if err := exec.Command("ip", "link", "delete", name, "type", "bridge").Run(); err != nil {
		return err
	}

	return nil
}
