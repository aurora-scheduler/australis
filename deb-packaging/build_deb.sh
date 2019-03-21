#!/bin/bash

docker build . -t australis_deb_builder

docker run --rm -v $HOME/go/pkg/mod:/go/pkg/mod -v $(pwd)/..:/australis australis_builder
