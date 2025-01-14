# Templates

The template is based on Cosmology's [chain-admin template](https://github.com/cosmology-tech/create-cosmos-app/tree/fd87039feee9568a86aa8d8d19edea8f4a78f599/templates/chain-template) at commit `fd87039feee9568a86aa8d8d19edea8f4a78f599`.

Some modifications have been done to make it work nicer with Ignite scaffolded chains.
An exhaustive list of changes compared to the original template can be found [here](./ignite-chain-admin.patch).

## Development Instructions

When upgrading the templates:

- checkout `github.com/cosmology-tech/create-cosmos-app` at the above mentioned commit.
- apply the git patch to the `chain-template` directory (`git apply ignite-chain-admin.patch`)
- merge upstream changes from main
- commit the changes (as a single commit, rewriting history if necessary -- `git reset $(git merge-base main $(git branch --show-current))`)
- export the changes to a patch file (`git diff main > ignite-chain-admin.patch`)
