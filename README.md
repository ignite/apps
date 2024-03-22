# **Ignite Apps**

Welcome to the Ignite App repository, a hub for enhancing blockchain app development with Ignite CLI. Our goal is to provide a platform where developers can find and share tools and insights, making blockchain application development more efficient and insightful.

## **About Ignite Apps**

Ignite Apps aims to extend the functionality of Ignite CLI, offering both official and community-contributed integrations. These integrations are designed to streamline development processes and offer valuable insights for blockchain app developers.

### **Official Ignite Apps**

- Developed by the core Ignite engineering team.
- Rigorously tested and fully supported.
- To submit your community-built app for official inclusion, please follow our submission guidelines.

### **Community Build Apps**

- A space for developers to share their Ignite integrations.
- Open for contributions via pull requests.
- To have your app reviewed for inclusion in the official section, please indicate this in your submission.

## **Getting Started**

### **How to Install an App**

```bash
ignite app install -g github.com/ignite/apps/[app-name]
```

Example: Installing the Hermes app

```bash
ignite app install -g github.com/ignite/apps/hermes
```

For more details, see [Installing Ignite Apps](https://docs.ignite.com/apps/using-apps).

### **How to Create an App**

Scaffold your Ignite app with one simple command:

```bash
ignite scaffold app path/to/your/app
```

Afterwards, install using:

```bash
ignite app install -g path/to/your/app
```

For more information, refer to [Creating Ignite Apps](https://docs.ignite.com/apps/developing-apps).

## **Contributing**

We welcome and appreciate new contributions. If you have an idea or an app that can benefit Ignite users, please follow our contribution guidelines.

- Fork the repository.
- Create your feature branch (**`git checkout -b feature/AmazingFeature`**).
- Commit your changes (**`git commit -am 'Add some AmazingFeature'`**).
- Push to the branch (**`git push origin feature/AmazingFeature`**).
- Open a pull request.

For detailed contribution guidelines, please refer to [CONTRIBUTING.md](CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## **Repository Structure**

Each directory in the root of this repository is a Go module containing an Ignite App package, with each app having its own go.mod file. This structure ensures modularity and ease of management for each app within the Ignite ecosystem.

## **Support and Feedback**

For support, questions, or feedback, please open an issue in this repository.

## **License**

This project is licensed under the [Copyright License](LICENSE) - see the [LICENSE](LICENSE) file for details.

## **Community**

- Join the community conversations on [Discord](https://discord.com/invite/ignite) or [X/Twitter](https://twitter.com/ignite).
- Follow the project's progress and updates.

## Developer instruction

- Clone this repo locally.
- Scaffold your app: `ignite app scaffold my-app`
- Add the folder to the `go.work`.
- Add your cobra commands into `debug/main.go` and the module replace to the `debug/go.mod` for a easy debug.
- Add the plugin: `ignite app add -g ($GOPATH)/src/github.com/ignite/apps/my-app`
- Test with Ignite.