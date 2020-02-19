# Australis

A light-weight client for [Aurora Scheduler](https://aurora-scheduler.github.io/) built using [gorealis](https://github.com/aurora-scheduler/gorealis).

## Usage
See the [documentation](docs/australis.md) for more information.

## Status
Australis is a work in progress and does not support all the features of Aurora Scheduler.

### Build locally
This project uses go mods. To build locally run:

`$ go build -o australis main.go`

### Building debian package
From the inside of the deb-packaging folder, run [build_deb.sh](deb-packaging/build_deb.sh)
