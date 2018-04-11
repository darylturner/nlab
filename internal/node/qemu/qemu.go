package qemu

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/darylturner/nlab/internal/config"
	"github.com/darylturner/nlab/internal/network"
	"github.com/darylturner/nlab/internal/network/netlinux"
	log "github.com/sirupsen/logrus"
)

func New(cfg config.NodeConf) *QemuNode {
	return &QemuNode{cfg}
}

type QemuNode struct {
	config.NodeConf
}

func (q QemuNode) Run(cfg *config.Topology, dryRun bool) error {
	var index int // need to find which node we are within the topology
	for i, nd := range cfg.Nodes {
		if nd.Tag == q.Tag {
			index = i
		}
	}

	serialPortBase := 50000 // need to make this more dynamic
	telnetPort := serialPortBase + index

	qemuArgs := []string{
		"-name", q.Tag, "-daemonize",
		"-smp", q.Resources.CPU,
		"-pidfile", fmt.Sprintf("/var/run/nlab/%v/%v.pid", cfg.Tag, q.Tag),
		"-m", q.Resources.Memory,
		"-display", "none", "-serial", fmt.Sprintf("telnet::%v,nowait,server", telnetPort),
	}

	for _, disk := range q.Resources.Disks {
		cmd := []string{"-drive", fmt.Sprintf("format=%v,file=%v", disk.Format, disk.File)}
		qemuArgs = append(qemuArgs, cmd...)
	}

	if q.Resources.CDROM != "" {
		qemuArgs = append(qemuArgs, "-cdrom "+q.Resources.CDROM)
	}

	virtIO := q.Network.VirtIO // virtio support specified?
	if q.Network.Management == true {
		tapName := netlinux.TapUID(cfg.ManagementBridge, q.Tag)
		qemuArgs = append(qemuArgs, linkCmd(cfg.ManagementBridge, tapName, virtIO)...)
	}
	for _, link := range q.Network.Links {
		tapName := netlinux.TapUID(link, q.Tag)
		qemuArgs = append(qemuArgs, linkCmd(link, tapName, virtIO)...)
	}

	if !dryRun {
		if out, err := exec.Command("kvm", qemuArgs...).CombinedOutput(); err != nil {
			return errors.New(fmt.Sprintf("%v: %v", err, string(out)))
		} else {
			log.WithFields(log.Fields{
				"tag":    q.Tag,
				"serial": fmt.Sprintf("telnet://localhost:%v", telnetPort),
			}).Info("running")
		}
	} else {
		fmt.Println("kvm " + strings.Join(qemuArgs, " "))
	}

	return nil
}

func (q QemuNode) Stop(cfg *config.Topology) error {
	pidPath := fmt.Sprintf("/var/run/nlab/%v/", cfg.Tag)

	pidBytes, err := ioutil.ReadFile(pidPath + fmt.Sprintf("%v.pid", q.Tag))
	if err != nil {
		return err
	}

	pidString := strings.TrimSpace(string(pidBytes))
	pid, err := strconv.Atoi(pidString)
	if err != nil {
		panic(err)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if err := proc.Kill(); err != nil {
		return err
	} else {
		log.WithFields(log.Fields{
			"tag": q.Tag,
		}).Info("node stopped")
	}

	return nil
}

func linkCmd(link, tap string, virtio bool) []string {
	drv := "e1000"
	if virtio {
		drv = "virtio-net-pci"
	}
	return []string{
		"-device", fmt.Sprintf("%v,netdev=%s,mac=%s", drv, link, network.RandomMAC()),
		"-netdev", fmt.Sprintf("tap,id=%s,ifname=%s,script=no", link, tap),
	}
}
