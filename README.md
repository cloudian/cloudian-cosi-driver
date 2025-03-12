# Cloudian COSI Driver

This repository contains the Cloudian [Container Object Storage Interface (COSI)](https://github.com/kubernetes-sigs/container-object-storage-interface) driver, compatible with a HyperStore backend.

# Getting Started

To deploy the pre-built containers, see [Quick Start](docs/quick-start.md).

# Developing

The repository includes a dev container, inside which it can be built and tested.

Building requires `golang` v1.23 or later, `docker`, and `make` to be installed.

### Building

To build the code and an up to date image:
```
make image
```

### Testing

To build a dev version of the image, then bring up a k3d test environment and load this image see `test/deploy.sh`.
