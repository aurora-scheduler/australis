#!/bin/bash

# Build debian package
cd /australis
debuild -d -us -uc -b

# Move resulting packages to the dist folder
mkdir -p /australis/dist
mv /australis_*_amd64* /australis/dist
