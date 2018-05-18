# nlab
CLI tool for managing Linux/KVM based network labs.

## Usage

```
# nlab --help
nlab can be used to create to simulate complicated network
topologies and launch KVM virtual-machines with sane defaults.
Tap/bridge and UDP pseudo-wire topologies supported.

Usage:
  nlab [command]

Available Commands:
  help        Help about any command
  network     Create and destroy Linux bridge/taps
  run         Run virtual machines
  stop        Stop virtual machines

Flags:
  -h, --help      help for nlab
  -j, --json      Output formatted as JSON to stdout
      --version   version for nlab

Use "nlab [command] --help" for more information about a command.
```

Configuration format is in JSON as a lingua franca. Config can be read
in through STDIN so anything that compiles to JSON should be fine to use.
```
# nlab network create lab.json
# nlab run -j lab.json
{"level":"info","msg":"running","serial":"telnet://localhost:50000","tag":"vmx0-re","time":"2018-05-04T11:43:11+01:00"}
{"level":"info","msg":"running","serial":"telnet://localhost:50001","tag":"vmx0-cp","time":"2018-05-04T11:43:11+01:00"}
{"level":"info","msg":"running","serial":"telnet://localhost:50002","tag":"vmx-peer1","time":"2018-05-04T11:43:11+01:00"}
# nlab network destroy lab.json
or
# convert2json lab.toml | nlab network create -
```

## Network Options

nlab was created specifically for running virtualised networks. The network topology is created simply by tagging machines with the required segments. The segments will be attached to the NICs
in the order they are specified in the array.

The default method of provisioning the underlying topology is Linux tap/bridge. Linux tap/bridge is high performance and allows good introspection using tcpdump on the associated bridge and tap adapters from the host. Some low level protocols may be intercepted by the host however and not forwarded as may be expected.

UDP pseudo-wires can also be used to provide better support for layer-2 protocols. LLDP, CDP and spanning-tree should all work correctly when using pseudo-wires. Ports will be assigned from 30000 automatically as required.

## Limitations

nlab does no bootstrapping of the nodes. A configuration however can be placed in a base image and cloned or copied. It shouldn't get in the way of any existing configuration management or ZTP tooling as network management can be provisioned onto a specified host bridge.

If using the management option only the first NIC is currently supported. I so far haven't found a network machine which doesn't use this convention.

## Recommendations

Using ZFS for the virtual machine disks works very well IMHO. Disk images can be written to ZFS sparse volumes; base images created using snapshots; and instantly cloned to form network hosts.
