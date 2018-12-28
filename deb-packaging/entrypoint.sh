#!/bin/bash


# Temporary fix for a go mods bug
rm /australis/go.sum

# Build debian package
cd /australis
debuild -d -us -uc -b

# Move resulting packages to the dist folder
mkdir -p /australis/dist
mv /australis_*_amd64* /australis/dist