# nlab
CLI tool for managing Linux/KVM based network labs.

## Usage

```
# nlab --help
nlab can be used to create Linux bridges and taps to
simulate complicated network topologies and launch KVM
virtual-machines with sane defaults.

Usage:
  nlab [command]

Available Commands:
  help        Help about any command
  network     Create and destroy Linux bridge/taps
  run         Run virtual machines

Flags:
  -h, --help   help for nlab

Use "nlab [command] --help" for more information about a command.
```

Configuration format is in JSON as a lingua franca. Config can be read
in through STDIN so anything that compiles to JSON should be fine to use.
```
# nlab network create lab.json
# nlab run lab.json
# nlab network destroy lab.json
or
# convert2json lab.toml | nlab create -
```
