# Explorer

Explorer is an Ignite App that enhances the [Ignite CLI](https://github.com/ignite/cli).

It integrates the TUI [Gex explorer](https://github.com/ignite/gex), which provides a real-time feed of your blockchain, enabling you to monitor activity and test your Ignite-based blockchains effectively.
Additionally, it integrates the [Ping.pub](https://ping.pub) explorer, which allows you to view your blockchain's transactions and blocks in a user-friendly web interface.

## Features

- **Real-time Blockchain Feed**: Get live updates and monitor blockchain events.
- **Seamless Integration**: Works effortlessly with blockchains scaffolded using Ignite.
- **Versatile Deployment**: Connect using various flexible testing and deployment methods.

## Prerequisites

Before using the Explorer app, ensure your environment meets the following requirements:

- **Ignite CLI**: Version `v28.3.0` or higher.

## Usage

Ensure your blockchain is running and the RPC server is accessible. You can start your blockchain using the following command:

```sh
ignite chain serve
```

### Ping.pub

To start the web explorer and connect it to your blockchain's RPC server, use the following command:

```sh
ignite explorer pingpub --port 8080
```

This command will start a web server on port `8080`, allowing you to access the Ping.pub explorer at `http://localhost:8080`.

### Gex

To start the TUI explorer and connect it to your blockchain's RPC server, use the following command:

```sh
ignite explorer gex --rpc-address http://localhost:26657
```

This command will display a live feed of blockchain activities.

For more information on how to use the TUI, check out the [Gex](https://github.com/ignite/gex) repository.
