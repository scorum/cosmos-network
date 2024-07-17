<div align="center">
  <h1> Scorum Network </h1>
</div>
<div align="center">
  <a href="https://github.com/scorum/cosmos-network/blob/main/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/scorum/cosmos-network.svg" />
  </a>
  <a href="https://docs.cosmos.network/main">
    <img alt="Cosmos Version: v0.46.11" src="https://img.shields.io/badge/cosmos_sdk-v0.46.11-blueviolet" />
  </a>

  <img alt="Mainnet Version: v1.0.3" src="https://img.shields.io/badge/mainnet-v1.0.3-green" />
  <img alt="Candidate Version: v1.1.0" src="https://img.shields.io/badge/candidate-v1.1.0-blue" />
</div>

# Introducing Scorum Cosmos Network

Scorum’s blockchain protocol is in the [cosmos](https://github.com/cosmos/cosmos-sdk) family.

### Projects currently working on Scorum Cosmos Network
- [Aviatrix](https://aviatrix.bet)

### Projects on which work is underway to transfer to the Scorum Cosmos Network
- [Blogging platform](https://scorum.com)
- [Betting platform](https://betscorum.com)

### Current and future functions of Scorum Cosmos Network
- Blogging platform where authors and readers will be rewarded for creating and engaging with content
- Commission-free betting exchange where fans can place bets against each other using Scorum Coins (SCR)
- Direct access to the Cosmos Network ecosystem
- Aviatrix NFT marketplace

### Public Announcement & Discussion
Find out more as we take the project public through the following channels:
- Get the latest updates and chat with us on [Telegram](https://telegram.me/SCORUM)

---

# Run Local Node Quick Start
This assumes that you're running Linux or MacOS and have installed [Go 1.19+](https://golang.org/dl/).  This guide helps you:

* build and install Scorum Network
* allow you to name your node
* add seeds to your config file
* download genesis state
* start your node
* use scorumd to check the status of your node.


If you already have a previous version of Scorum Network installed:
```
rm -rf ~/.scorum
```

## Development

### Requirements
To build project you should have:
- go >= 1.19
- protobuf
- buf
- golangci-lint

### Guide

#### Build and install
Build and install are available with `make`
```shell
make build
```

or 

```shell
make install
```

#### Install tools
```shell
make install-proto
make install-linter
```

#### Generate proto
```shell
make generate-proto
```

#### Generate swagger
To generate swagger you need:
- buf
- cloned to gopath github.com/cosmos/cosmos-sdk with correct version
- go-swagger-merger
```shell
make generate-proto-swagger
```

#### Run tests
```shell
make test
```

Also you can run non-determinism simulation
```shell
make test-determinism
```

#### Run linter
```shell
make lint
```

#### Start local node
Node's home is a directory test in the root of the repo. There will be available `test` account in the test directory. 

```shell
make local-init
make local-start
```

To reset state
```shell
make local-reset
```

#### Run cli with the local node
Be sure you installed `scorumd` or use the built one from `./build`

Retrieve info about test account
```shell
scorumd keys list  --keyring-backend test --keyring-dir test
```

Broadcast transaction
```shell
scorumd tx aviatrix create-plane 06c444a4-1ca4-11ee-be56-0242ac120003 scorum1z2ptzs6a2f8dfp4et22x36a53x0w6cma6pplca my-plane black --keyring-backend test --keyring-dir test --broadcast-mode block --from test
```

Query command
```shell
scorumd query bank balances scorum1qujv377hmegm46vnllm20y9c56m5nztserc07a
```

---

# No Support & No Warranty

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
IN THE SOFTWARE.

# License

© 2023 Scorum. All right reserved.

Licensed under the Apache v2 License.
