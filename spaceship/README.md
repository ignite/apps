# Spaceship

Spaceship is an Ignite App designed to extend the [Ignite CLI](https://github.com/ignite/cli) by providing tools to deploy blockchain applications via SSH.

## Prerequisites

* **Ignite CLI**: Version `v28.4.0` or higher is required.
* **Blockchain Scaffold**: A blockchain scaffolded using Ignite.

## Usage

Spaceship provides multiple ways to connect to your SSH server for deployment:

```sh
ignite spaceship deploy root@127.0.0.1 --key $HOME/.ssh/id_rsa
ignite spaceship deploy 127.0.0.1 --user root --key $HOME/.ssh/id_rsa
ignite spaceship deploy 127.0.0.1 --user root --password password
ignite spaceship deploy root@127.0.0.1 --key $HOME/.ssh/id_rsa --key-password key_password
```

Each command initiates a build of the blockchain binary and sets up the chain's home directory based on the configuration. The app then connects to the specified SSH server, establishes workspaces, transfers the binary, and executes it using a runner script.

The workspaces are organized under `$HOME/workspace/<chain-id>` and include:
- **Binary Directory**: `$HOME/workspace/<chain-id>/bin` - Contains the chain binary.
- **Home Directory**: `$HOME/workspace/<chain-id>/home` - Stores chain data.
- **Log Directory**: `$HOME/workspace/<chain-id>/log` - Holds logs of the running chain.
- **Runner Script**: `$HOME/workspace/<chain-id>/run.sh` - A script to start the binary in the background using `nohup`.
- **PID File**: `$HOME/workspace/<chain-id>/spaceship.pid` - Stores the PID of the currently running chain instance.

### Managing the Chain

To manage your blockchain deployment, use the following commands:

- **Check status**:

```sh
ignite spaceship status root@<ip-address> --key $HOME/.ssh/id_rsa
```

- **View logs**:

```sh
ignite spaceship log root@<ip-address> --key $HOME/.ssh/id_rsa
```

- **Watch the logs in real time**:

```sh
ignite spaceship log root@<ip-address> --key $HOME/.ssh/id_rsa --real-time
```

- **Restart the chain**:

```sh
ignite spaceship restart root@<ip-address> --key $HOME/.ssh/id_rsa
```

- **Stop the chain**:

```sh
ignite spaceship stop root@<ip-address> --key $HOME/.ssh/id_rsa
```

To redeploy the chain on the same server without overwriting the home directory, use the `--init-chain` flag to reinitialize the chain if necessary.

## Faucet

You can deploy your chain along with a faucet application by passing the faucet flag to the deploy command:

```sh
ignite spaceship deploy root@127.0.0.1 --key $HOME/.ssh/id_rsa --faucet
```

You can also specify the faucet port:

```sh
ignite spaceship deploy root@127.0.0.1 --key $HOME/.ssh/id_rsa --faucet --faucet-port 8077
```

To check the faucet logs, pass the parameter `faucet` to the `--app` flag in the log command:

```sh
ignite spaceship log root@127.0.0.1 --key $HOME/.ssh/id_rsa --app faucet
```

### Config

You can override the default [chain configuration](https://docs.ignite.com/references/config#validators) by using the Ignite configuration file. Validators' node configuration files are stored in the data directory. By default, Spaceship initializes the chain locally in a temporary folder using the Ignite config file and then copies the configuration to the remote machine at `$HOME/workspace/<chain-id>/home`.

Configuration resets are performed by Ignite when necessary, especially when using the `--init-chain` flag or if the chain was not previously initialized.

**Example Ignite config**:

```yaml
validators:
  - name: alice
    bonded: '100000000stake'
    app:
      pruning: "nothing"
    config:
      moniker: "mychain"
    client:
      output: "json"
```

For more information, please refer to the [Ignite documentation](https://docs.ignite.com).
