# Ignite App Registry

## Overview

The Ignite App Registry is a directory designed to catalog applications built within the Ignite Apps ecosystem. It aims to facilitate discovery, collaboration, and information sharing among developers and users. By submitting your app to the registry, you make it easier for the community to learn about and engage with your project.

## How to Add Your Ignite App

To add your app to the Ignite App Registry, follow these steps:

1. **Fork the Repository**: Start by forking the `ignite/apps` repository to your own GitHub account.
2. **Set Up Your App Directory**: Inside the repository and inside `_registry` directory, create a file named after your username, repository name and app name: `username.repository.app_name.json`.
3. **Configure Your App**: The previously added file should contain all the relevant details about your app according to the template provided below.
4. **Submit a Pull Request**: Once the above steps are completed, submit a pull request to the original `ignite/apps` repository. Your pull request will be reviewed by the maintainers, and once approved, your app will be listed in the registry.

## `app.json` File Structure**

Below is the template for the `username.repository.app_name.json` file. Replace each placeholder with the specific details of your app.

```json
{
    "appName": "Your App Name",
    "appDescription": "A brief description of what your app does.",
    "ignite": ">28.3.0",
    "dependencies": {
      "docker": ">23.0.5"
    },
    "cosmosSDK": ">0.50.4",
    "features": [
      "Feature 1",
      "Feature 2",
      "Additional features..."
    ],
    "wasm": false,
    "authors": [
      {
        "name": "Author Name",
        "email": "email@example.com",
        "website": "Optional author or company website"
      }
    ],
    "repository": {
      "url": "URL to the app's repository"
    },
    "documentationUrl": "URL to the app's documentation",
    "license": {
      "name": "MIT",
      "url": "github.com/username/app/LICENSE.md"
    },
    "keywords": ["keyword1", "keyword2", "Useful for search and categorization"],
    "supportedPlatforms": ["mac", "linux"],
    "socialMedia": {
      "x": "Optional X handle",
      "telegram": "Optional Telegram group",
      "discord": "Optional Discord server",
      "reddit": "Optional Reddit page",
      "website": ""
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

## Best Practices

- **Keep Information Up-to-Date**: Regularly update your app json file to reflect the latest releases, compatibility changes, and other relevant information.
- **Detailed Descriptions**: Provide clear and concise descriptions of your app's features and functionality.
- **Engage with the Community**: Participate in discussions, address issues, and provide support to encourage engagement with your app.

## Support and Contributions

The Ignite App Registry is a community-driven project. We welcome contributions, suggestions, and feedback. If you encounter issues or have ideas for improving the registry, please open an issue or submit a pull request.

## Disclaimer

Inclusion in the registry is not an official endorsement or recognition.
