# autodock

[![Build Status](https://travis-ci.org/prologic/autodock.svg)](https://travis-ci.org/prologic/autodock)
[![CodeCov](https://codecov.io/gh/prologic/autodock/branch/master/graph/badge.svg)](https://codecov.io/gh/prologic/autodock)
[![Go Report Card](https://goreportcard.com/badge/github.com/prologic/autodock)](https://goreportcard.com/report/github.com/prologic/autodock)
[![Image Layers](https://badge.imagelayers.io/prologic/autodock:latest.svg)](https://imagelayers.io/?images=prologic/autodock:latest)
[![GoDoc](https://godoc.org/github.com/prologic/autodock?status.svg)](https://godoc.org/github.com/prologic/autodock)

[autodock](https://github.com/prologic/autodock) is a Daemon for
Docker Automation which enables you to maintain and automate your Docker
infrastructure by reacting to Docker or Docker Swarm events.

## Supported plugins:

autodock comes with a number of plugins where each piece of functionality is
rovided by a separate plugin. Each plugin is "linked" to autodock to receive
Docker events and issue new Docker API commands.

The following list is a list of the currently available plugins:

- [autodock-cron](https://github.com/prologic/autodock)
  Provides a *Cron* like scheduler for Containers/Services
- [autodock-logger](https://github.com/prologic/autodock-logger)
  Logs Dockers Events

## Installation

### Docker

```#!bash
$ docker pull prologic/autodock
```

### Source

```#!bash
$ go install github.com/prologic/autodock
```

## Usage

### Docker

```#!bash
$ docker run -d -p 8000:8000 -v /var/run/docker.sock:/var/run/docker.sock prologic/autodock
```

### Source

```#!bash
$ autodock
```

## License

MIT
