# Explorer

Explorer is an Ignite App that enhances the [Ignite CLI](https://github.com/ignite/cli) by integrating
the [Gex explorer](https://github.com/ignite/gex). This app provides a real-time feed of your blockchain,
enabling you to monitor activity and test your Ignite-based blockchains effectively.

## Features

- **Real-time Blockchain Feed**: Get live updates and monitor blockchain events.
- **Seamless Integration**: Works effortlessly with blockchains scaffolded using Ignite.
- **Versatile Deployment**: Connect using various flexible testing and deployment methods.

## Prerequisites

Before using the Explorer app, ensure your environment meets the following requirements:

- **Ignite CLI**: Version `v28.3.0` or higher.
- Additionally, your blockchain should be scaffolded using the Ignite CLI.

## Usage

To start the Explorer and connect it to your blockchain's RPC server, use the following command:

```sh
ignite explorer gex --rpc-address http://localhost:26657
```

This command will display a live feed of blockchain activities.

## Additional Options

For more information, please refer to the [Ignite documentation](https://docs.ignite.com) or the [Gex](https://github.com/ignite/gex) repository.