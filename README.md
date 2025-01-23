# Ignite Apps

Welcome to the Ignite App repository, a hub for enhancing blockchain app development with Ignite CLI. Our goal is to provide a platform where developers can find and share tools and insights, making blockchain application development more efficient and insightful.

## About Ignite Apps

Ignite Apps aims to extend the functionality of Ignite CLI, offering both official and community-contributed integrations. These integrations are designed to streamline development processes and offer valuable insights for blockchain app developers.

### Official Ignite Apps

- Developed by the core Ignite engineering team.
- Rigorously tested and fully supported.
- To submit your community-built app for official inclusion, please follow our submission guidelines.

## Getting Started

### How to Install an App

```bash
ignite app install -g github.com/ignite/apps/[app-name]
```

Example: Installing the Hermes app

```bash
ignite app install -g github.com/ignite/apps/hermes
```

For more details, see [Installing Ignite Apps](https://docs.ignite.com/apps/using-apps).

### How to Create an App

Scaffold your Ignite app with one simple command:

```bash
ignite scaffold app path/to/your/app
```

Afterwards, install using:

```bash
ignite app install -g path/to/your/app
```

For more information, refer to [Creating Ignite Apps](https://docs.ignite.com/apps/developing-apps).

## Contributing

We welcome and appreciate new contributions. If you have an idea or an app that can benefit Ignite users, please follow our contribution guidelines.

- Fork the repository.
- Create your feature branch (**`git checkout -b feature/AmazingFeature`**).
- Commit your changes (**`git commit -am 'Add some AmazingFeature'`**).
- Push to the branch (**`git push origin feature/AmazingFeature`**).
- Open a pull request.

For detailed contribution guidelines, please refer to [CONTRIBUTING.md](CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## Repository Structure

Each directory in the root of this repository is a Go module containing an Ignite App package, with each app having its own go.mod file.
This structure ensures modularity and ease of management for each app within the Ignite ecosystem.

app-name/
├── cmd/
│   └── command_one.go
│   └── command_two.go
├── integration/
│   └── app_test.go
├── changelog.md
├── go.mod
├── go.sum
├── main.go
└── README.md

The actual implementation of the app is in the root directory, while the `cmd` directory contains the commands that are exposed to the user.

## Support and Feedback

For support, questions, or feedback, please open an issue in this repository.

### Community Build Apps

The Ignite Apps ecosystem thrives on community contributions, offering a space for developers to share, showcase, and collaborate on Ignite integrations. If you've built an app that leverages Ignite CLI and wish to share it with the community, we encourage you to do so by following these guidelines:

### Create Your Repository

1. **Setup Your Project:** Begin by creating a new repository for your project on GitHub. This repository will host your app's code, documentation, and other relevant materials.
2. **Choose a License:** Choosing an appropriate license for your project is essential. This ensures that others know how they can legally use, modify, and distribute your app. GitHub's licensing guide can help you decide which license to use if you're unsure.
3. **Prepare Your Documentation:** Your repository should include clear and concise documentation that covers:
    - An overview of your app and its features.
    - Installation instructions, including any prerequisites or dependencies.
    - Usage examples to help users get started quickly.
    - Contribution guidelines if you're open to receiving contributions from others.

### Integrate with the Ignite App Registry

Once your repository is set up and ready to be shared, we invite you to add your app to the Ignite App Registry. This registry is a curated list of community-built apps, making it easier for users to discover and utilize your contributions.

1. **Fork the App-Registry Repository:** Visit the Ignite App-Registry repository and fork it to your GitHub account.
2. **Add Your App:** In your fork, add a new entry for your app in the registry. This should include the name of your app, a brief description, and a link to your app's repository.
3. **Submit a Pull Request:** Once you've added your app to the registry, submit a pull request to the original App-Registry repository. Our team will review your submission and, if everything is in order, merge it into the registry.

### Engage with the Community

Sharing your app is just the beginning. Engage with the Ignite community by:

- Responding to issues and pull requests in your repository.
- Promoting your app on social media and Ignite community channels.
- Keeping your app updated and improving it based on user feedback.

We're excited to see your contributions and how they will enrich the Ignite Apps ecosystem. Together, we can build a diverse and vibrant community of blockchain application developers.

## License

This project is licensed under the [Copyright License](LICENSE) - see the [LICENSE](LICENSE) file for details.

## Community

- Join the community conversations on [Discord](https://discord.com/invite/ignite) or [X/Twitter](https://twitter.com/ignite).
- Follow the project's progress and updates.

## Developer instruction

- Clone this repo locally.
- Scaffold your app: `ignite app scaffold my-app`
- Add the folder to the `go.work`.
- Add your cobra commands into `debug/main.go` and the module replace to the `debug/go.mod` for a easy debug.
- Add the plugin: `ignite app add -g ($GOPATH)/src/github.com/ignite/apps/my-app`
- Test with Ignite.
