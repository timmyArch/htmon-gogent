# htmon-gogent
Http monitoring agent written in Go.

## Installation

* Install go package

```bash
  GOPATH=...
  PATH=$PATH:$GOPATH/bin
  go get github.com/timmyArch/htmon-gogent
```

* After installtion you should be able to run

```bash
  htmon-gogent --help
```

## Usage

* Create configuration

```yaml
server:
  user: <User>
  password: <Password>
  apiurl: https://htmon.moo.gl/api/v1
checks:
  - check:
      interval: 2
      expire: 60
      metric: 'process::$$placeholder$$'
      command: 'pgrep -fla $$placeholder$$ | grep -v $$placeholder$$'
		# uses variable from schema
    placeholder: htmon_processes
  - check:
      interval: 2
      expire: 60
      metric: 'process::$$placeholder$$'
      command: 'echo $$placeholder$$'
		# iterates over placeholders
    placeholder: [ 'nginx', 'grafana', 'test', 'foo', 'moo' ]
  - check:
      interval: 2
      expire: 60
      metric: 'metric::moo'
      command: 'echo test'
```

* Try a test run

```bash
	htmon-gogent agent --config <Path> --test
```

* Start agent loop

```bash
	htmon-gogent agent --config <Path>
```

* Print loaded schema with spoofed hostname.

```bash
	htmon-gogent agent --config <Path> --schema --spoof-hostname=example.org
```

