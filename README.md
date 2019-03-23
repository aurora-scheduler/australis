# Australis

A light-weight client for [Apache Aurora](https://aurora.apache.org/) built using [gorealis](https://github.com/paypal/gorealis).

## Usage
See the [documentation](docs/australis.md) for more information.

## Status
Australis is a work in progress and does not support all the features of Apache Aurora.

### Build locally
This project uses go mods. To build locally run:

`$ go build -o australis main.go`

### Building debian package
From the inside of the deb-packaging folder, run [build_deb.sh](deb-packaging/build_deb.sh)