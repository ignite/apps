# CosmWasm Plugin for Ignite CLI

**This repo contains** an Ignite CLI plugin that can help you to add CosmWasm module support for newly scaffolded app. 

## Get started
First scaffold a chain with [Ignite CLI](https://docs.ignite.com).

1. Install Ignite CLI:
```
$ curl https://get.ignite.com/cli! | bash
```
2. Scaffold a chain
```
$ ignite scaffold chain planet
$ cd planet
```

3. Clone this repo locally

4. Add cosmwasm plugin within your project directory:
```
ignite plugin add /absolute/path/to/plugin/cosmwasm-plugin
```

or globally

```
ignite plugin add -g /absolute/path/to/plugin/osmwasm-plugin
```

5. Run command
```
$ ignite cosmwasm-plugin add
```

6. Launch your chain with CosmWasm support

```
$ ignite chain serve
```
`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.


## Testing Smart-Contracts

Testing CosmWasm smart-contract capabilities in your chain, can be performed by following these steps:

0. Prerequisites: install [rust](https://www.rust-lang.org/tools/install) and [docker](https://www.docker.com/) (for smart-contact optimization)

1. Launch chain locally:
```
$ ignite chain serve
```

2. Prepare cosmwasm smart-contract for deployment
```
$ git clone https://github.com/CosmWasm/cw-examples
$ cd cw-examples
$ cd contracts/nameservice
$ RUSTFLAGS='-C link-arg=-s' cargo wasm
```

Optimize contract to reduce gas
```
$ docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/rust-optimizer:0.12.12
```

For ARM64 machines:
```
$ docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/rust-optimizer-arm64:0.12.12
```

3. Deploy smart-contract to your chain
```
$ planetd tx wasm store artifacts/cw_nameservice.wasm --from alice --gas=10000000 
```

4. Check your smart-contract on chain
```
$ planetd query wasm list-code
code_infos:
- code_id: "1"
  creator: cosmos1a4w3vwxa8zhy57nyy9kw7zhdasr4ue4wld0zpn
  data_hash: DF3C9BC1341322810523AABCA28CC3FCDCA021C85061743967CE3D20F5580093
pagination: {}

# download wasm and diff with origin
$ CODE_ID=$(planetd query wasm list-code --output json | jq -r '.code_infos[0].code_id')
$ planetd query wasm code $CODE_ID download.wasm
$ diff artifacts/cw_nameservice.wasm download.wasm
```

5. Instantiate smart-contract, in CosmWasm deployment and instantiation are 2 different steps

```
$ INIT='{"purchase_price":{"amount":"100","denom":"stake"},"transfer_price":{"amount":"999","denom":"stake"}}'
$ planetd tx wasm instantiate 1 "$INIT" --from alice --chain-id "planet" --label "awesome name service" --no-admin

$ CONTRACT=$(planetd query wasm list-contract-by-code $CODE_ID --output json | jq -r '.contracts[-1]')
# check contract state (you should get contract address)
$ planetd query wasm contract $CONTRACT
address: cosmos14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s0phg4d
contract_info:
  code_id: "1"
  creator: cosmos1a4w3vwxa8zhy57nyy9kw7zhdasr4ue4wld0zpn
  label: awesome name service
```

6. Interact with smart-contract
```
# purchase name
$ planetd tx wasm execute $CONTRACT '{"register":{"name":"test"}}' --amount 100stake --from alice $TXFLAG -y

# query registered name (you should see your address as owner of the name)
$ NAME_QUERY='{"resolve_record": {"name": "test"}}'
$ planetd query wasm contract-state smart $CONTRACT "$NAME_QUERY" --output json
```

## Compatibility

| Ignite CLI  | Cosmos SDK  | IBC       | wasmd                                                         |
|-------------|-------------|-----------|---------------------------------------------------------------|
| v0.27.1     | v0.47.4     | v7.2.0    | v0.41.0                                                       |

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)
- [CosmWasm](https://cosmwasm.com/)
- [Wasm module](https://github.com/CosmWasm/wasmd)
- [CosmWasm Smart Contract Tutorial](https://medium.com/haderech-dev/smart-contract-tutorial-3-cosmwasm-805860c91a88)