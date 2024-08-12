# Release Process

This document outlines the release process for Scorum Network.

## Release Procedure

Once the Scorum team decide to release new version the release candidate tag is created from the main branch.
It goes through internal testing process until it is ready to be released on the mainnet.  

Once the release is prepared the team creates the final tag to upgrade the network to.
The team broadcasts a proposal and announces about that on community channels.

## Upgrading

All upgrades support [cosmovisor](https://docs.cosmos.network/main/build/tooling/cosmovisor). If a validator doesn't use cosmovisor the sources can be built manually by
```
git clone git@github.com:scorum/cosmos-network.git
git checkout <tag>
make install
```

*Notice*. `go` and `make` must be installed.

## Changelog

All changes across releases are described in the [changelog](.CHANGELOG.md).