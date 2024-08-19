# Build local image
```
make local-image
```

# Prepare networks

Change {index} to 0 or 1, depending it is 1st or 2nd node installation

0. Add chain-id param to init step in make file
```
$(V)$(SCORUMD) init --default-denom nscr local-network --chain-id scorum-{index}
```

1. Prepare fresh network genesis and config
```
make local-init
```

2. Rename test to test-local-{index}
```
mv test test-local-{index}
```

3. Start network
```
docker run -d --name test-local-{index} --rm -v ./test-local-{index}:/root/.scorum  scorumd-local scorumd start --api.enable true --grpc.enable true --rpc.laddr 'tcp://0.0.0.0:26657' --grpc.address '0.0.0.0:9090'
```

4. Get container address by
```
docker network inspect bridge
```

5. Start container to interact with network
```
docker run --rm -it -v ./test-local-{index}:/root/.scorum  scorumd-local /bin/sh
```

6. Create key and save mnemonic somewhere
```
scorumd keys add hermes --keyring-backend test
```

7. Mint gas to this address
```
scorumd tx scorum admin mint-gas --node 'tcp://{node-ip}:26657' --from test --keyring-backend test {hermes-key-address} 10000000000000000
```

Now we can consider network prepared


# Prepare hermes

1. Create config dir
```
mkdir test-hermes
```

2. Create hermes config file
```
nano test-hermes/config.toml
```

Put config and save it by ctrl+x.

*NOTE* Copy paste [[chains]] template. There are should be 2 [[chains]] blocks with id 'scorum-0' and 'scorum-1'
```
[global]
log_level = 'info'

[mode]

[mode.clients]
enabled = true
refresh = true
misbehaviour = true

[mode.connections]
enabled = true

[mode.channels]
enabled = true

[mode.packets]
enabled = true
clear_interval = 100
clear_on_start = true
tx_confirmation = true

[telemetry]
enabled = true
host = '127.0.0.1'
port = 3001

[[chains]]
id = 'scorum-{index}'
type = "CosmosSdk"
rpc_addr = 'http://{node-ip}:26657'
grpc_addr = 'http://{node-ip}:9090'
event_source = { mode = 'push', url = 'ws://{node-ip}:26657/websocket', batch_delay = '200ms' }
rpc_timeout = '15s'
trusted_node = true
account_prefix = 'scorum'
key_name = 'hermes'
store_prefix = 'ibc'
gas_price = { price = 1, denom = 'gas' }
gas_multiplier = 1.2
default_gas = 1000000
max_gas = 10000000
max_msg_num = 30
max_tx_size = 2097152
clock_drift = '5s'
max_block_time = '30s'
trusting_period = '14days'
trust_threshold = { numerator = '2', denominator = '3' }
```

3. Start hermes container
```
docker run --rm -it -v ./test-hermes:/home/hermes/.hermes --entrypoint "/bin/bash" informalsystems/hermes:1.10.2
```

4. Add keys to hermes
```
echo "{mnemonic}" | hermes keys add --chain scorum-{index} --mnemonic-file /dev/stdin
```

5. Create channel
```
hermes create channel --a-chain scorum-0 --b-chain scorum-1 --a-port transfer --b-port transfer --new-client-connection
```

6. Start relaying
```
hermes start
```

# Test ibc
1. Run ibc transfer from *test-local-0* terminal
```
scorumd tx ibc-transfer transfer --node 'tcp://{node-ip}:26657' --from test --keyring-backend test transfer channel-0 scorum17sxdedhm47y5s9w55ng0nflyqn0w5ffxwz3cuq 1000nscr
```

2. Query balance on *test-local-1* terminal
```
scorumd --node 'tcp://{node-ip}:26657' query bank balances scorum17sxdedhm47y5s9w55ng0nflyqn0w5ffxwz3cuq
```

You must see something like
```
- amount: "1000"
  denom: ibc/8529F970ED97E874805713A099E86A520D418AF2AC44335C7E1FED9D545C4020
```