# Ignite Apps

Official [Ignite CLI](https://ignite.com/cli) apps repository.

Each directory in the root of this repo must be a Go module containing
an Ignite App package, each one with its own `go.mod` file.


## Wasmd
In this repository we have add the required steps to add wasmd into your chain.
you can start by : 

``
    ignite wasmd init
``

currently it is a minimal version just adding the wiring but we are going to develop some other commands for smart contracts deployment and testing. 

In this regard, we have used a combination of placeholders and replacements.

here are already used placeholders in ignite that you can find in yaml file under changes section. and there were few places that we couldn't find any placeholder to fit and we used replacement for those parts and you can find it under replace section in yaml file.

### TODO:
- [ ] add version compatibility checking for cosmos and wasmd
- [ ] add smart contract features (new/test/deploy)
- [ ] IBC options and other parameters

