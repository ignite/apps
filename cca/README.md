# CCA App

This Ignite App is aimed to extend [Ignite CLI](https://github.com/ignite/cli) and bootstrap the development of a chain frontend.
It uses the widely used [Cosmology create-cosmos-app](https://github.com/cosmology-tech/create-cosmos-app) template and its libraries.

## Prerequisites

* Ignite v28.7.0 or later
* Node.js
* [Corepack](https://yarnpkg.com/corepack)

## Installation

```shell
ignite app install -g github.com/ignite/apps/cca
```

### Usage

Run the app using `ignite chain serve` command.
In another terminal, run the frontend using the following commands:

```shell
ignite s cca
cd web
yarn install
yarn dev
```

Learn more about Cosmos-Kit and Ignite in their respective documentation:

* <https://docs.ignite.com>
* <https://github.com/cosmology-tech/create-cosmos-app>
