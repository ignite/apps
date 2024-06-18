# Contributing to Ignite Apps

We warmly welcome contributions to the Ignite Apps project! This document provides guidelines for contributing to the Ignite Apps repository. By following these guidelines, you can help us maintain a healthy and sustainable open-source ecosystem.

## Code of Conduct

This project and everyone participating in it is governed by the [Ignite Apps Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the team.

## How to Contribute

### Reporting Bugs

Before submitting a bug report, please check if the issue has already been reported. If it has not, create a new issue and include:

- A clear title and description.
- As much relevant information as possible (e.g., version, environment).
- A code sample or an executable test case demonstrating the expected behavior that is not occurring.

### Suggesting Enhancements

Enhancement suggestions are also welcome. Open an issue and include:

- A clear title and description of the enhancement.
- Explain why this enhancement would be useful.
- Provide a use case or code examples if possible.

### Pull Requests

1. Fork the repository and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. Ensure your code adheres to the existing style in the project to maintain consistency.
4. Write clear, meaningful commit messages.
5. Include appropriate tests.
6. Add or update the documentation as necessary.
7. After submitting your pull request, verify that all status checks are passing.

### Development Setup

For setting up your local development environment, follow these steps:

```bash
# Clone the Apps directory
git clone https://github.com/ignite/apps.git && cd apps
# Scaffold a template for your own app
ignite scaffold app <name> && cd <name>
# Install your app
ignite app install -g $(pwd)
```

### Contribution Prerequisites

- Familiarity with Git and GitHub.
- Understanding of the project's technology stack and goals.

## Community

- Join the community conversations on [Discord](https://discord.com/invite/ignite) or [X/Twitter](https://twitter.com/ignite).
- Follow the project's progress and updates.

## Pull Request Process*

- Update the README.md or documentation with details of changes to the interface, if applicable.
- Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent.
- The pull request will be merged once it's reviewed and approved by the maintainers.

## **License**

See [LICENSE](LICENSE) for more information.

## Questions?

If you have any questions about contributing, please feel free to contact us.

Thank you for your interest in contributing to Ignite Apps!
