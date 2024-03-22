# **Ignite App Registry**

## **Overview**

The Ignite App Registry is a directory designed to catalog applications built within the Ignite Apps ecosystem. It aims to facilitate discovery, collaboration, and information sharing among developers and users. By submitting your app to the registry, you make it easier for the community to learn about and engage with your project.

## **How to Add Your Ignite App**

To add your app to the Ignite App Registry, follow these steps:

1. **Fork the Repository**: Start by forking the `ignite/apps` repository to your own GitHub account.
2. **Create a New Directory within `_ignite-app-registy`: Inside the repository and inside the directory `_ignite_app_registry`, create a new directory named after your app. Ensure the name is unique and descriptive.
3. **Add Your `app.json` File**: In your app's directory, add an `app.json` file. This file should contain all the relevant details about your app according to the template provided below.
4. **Submit a Pull Request**: Once you've added your app's directory and `app.json` file, submit a pull request to the original `ignite/apps` repository. Your pull request will be reviewed by the  maintainers, and once approved, your app will be listed in the registry.

## `app.json` File Structure**

Below is the template for the `app.json` file. Replace each placeholder with the specific details of your app.

```json
{
    "appName": "Your App Name",
    "appDescription": "A brief description of what your app does.",
    "latestRelease": {
      "version": "x.y.z",
      "date": "YYYY-MM-DD",
      "releaseNotesUrl": "URL to the release notes or changelog"
    },
    "igniteCompatibility": {
      "lastTestedVersion": "ignite vX.Y.Z"
    },
    "cosmosSdkCompatibility": {
      "lastWorkedVersion": "vX.Y.Z"
    },
    "features": [
      "Feature 1",
      "Feature 2",
      "Additional features..."
    ],
    "wasm": "Indicate if your app uses wasm (yes/no)",
    "authors": [
      {
        "name": "Author Name",
        "email": "email@example.com",
        "website": "Optional author or company website"
      }
    ],
    "repository": {
      "url": "URL to the app's repository",
      "type": "e.g., GitHub, GitLab"
    },
    "documentationUrl": "URL to the app's documentation",
    "license": "License type, e.g., MIT, Apache 2.0",
    "keywords": ["keyword1", "keyword2", "Useful for search and categorization"],
    "supportedPlatforms": ["platform1", "platform2", "e.g., osmosis, cosmoshub"],
    "socialMedia": {
      "twitter": "Optional Twitter handle",
      "telegram": "Optional Telegram group",
      "discord": "Optional Discord server",
      "reddit": "Optional Reddit page",
      "website": "Optional website page"
    },
    "donations": {
      "cryptoAddresses": {
        "cosmos": "cosmos1...",
        "otherSupportedCryptos": "address"
      },
      "fiatDonationLinks": "URL to Patreon, Ko-fi, etc."
    }
}

```

## **Best Practices**

- **Keep Information Up-to-Date**: Regularly update your `app.json` file to reflect the latest releases, compatibility changes, and other relevant information.
- **Detailed Descriptions**: Provide clear and concise descriptions of your app's features and functionality.
- **Engage with the Community**: Participate in discussions, address issues, and provide support to encourage engagement with your app.

## **Support and Contributions**

The Ignite App Registry is a community-driven project. We welcome contributions, suggestions, and feedback. If you encounter issues or have ideas for improving the registry, please open an issue or submit a pull request.

## Disclaimer

Inclusion in the registry is not an official endorsement or recognition.