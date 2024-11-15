export type ProductCategory =
  | 'cosmwasm'
  | 'cosmos-sdk'
  | 'frontend'

export type Product = {
  name: string;
  description: string;
  link: string;
  category: ProductCategory;
};

export const products: Product[] = [
  {
    name: 'Ignite',
    description:
      'Ignite makes developing, growing, and launching blockchain projects faster than ever before.',
    link: 'https://ignite.com',
    category: 'cosmos-sdk',
  },
  {
    name: 'Ignite Tutorials',
    description:
      'Learn How to Build Cutting-Edge Blockchains.',
    link: 'https://tutorials.ignite.com',
    category: 'cosmos-sdk',
  },
  {
    name: 'Cosmos SDK',
    description:
      'Get a quick introduction to the Cosmos SDK and its key features, including its modular architecture and developer-friendly tools.',
    link: 'https://docs.cosmos.network',
    category: 'cosmos-sdk',
  },
  {
    name: 'Cosmos Kit',
    description:
      'A wallet adapter for react with mobile WalletConnect support for the Cosmos ecosystem.',
    link: 'https://cosmology.zone/products/cosmos-kit',
    category: 'frontend',
  },
  {
    name: 'Ignite Videos',
    description:
      'Learn about Ignite development with video tutorials.',
    link: 'https://youtube.com/@ignitehq',
    category: 'frontend',
  },
  {
    name: 'Telescope',
    description:
      'A TypeScript Transpiler for Cosmos Protobufs to generate libraries for Cosmos blockchains.',
    link: 'https://cosmology.zone/products/telescope',
    category: 'cosmos-sdk',
  },
  {
    name: 'Interchain UI',
    description:
      'A simple, modular and cross-framework component library for Cosmos ecosystem.',
    link: 'https://cosmology.zone/products/interchain-ui',
    category: 'frontend',
  },
  {
    name: 'Chain Registry',
    description:
      'Get chain and asset list information from the npm package for the Official Cosmos chain registry.',
    link: 'https://cosmology.zone/products/chain-registry',
    category: 'frontend',
  },
  {
    name: 'Create Cosmos App',
    description:
      'One-Command Setup for Modern Cosmos dApps. Speed up your development and bootstrap new web3 dApps quickly.',
    link: 'https://cosmology.zone/products/create-cosmos-app',
    category: 'frontend',
  },
  {
    name: 'CosmWasm Academy',
    description:
      'Master CosmWasm and build your secure, multi-chain dApp on any CosmWasm chain!',
    link: 'https://cosmology.zone/learn/ts-codegen',
    category: 'cosmwasm',
  },
  {
    name: 'Cosmology Videos',
    description:
      'How-to videos from the official Cosmology website, with learning resources for building in Cosmos.',
    link: 'https://cosmology.zone/learn',
    category: 'frontend',
  },
  {
    name: 'Next.js',
    description: 'A React Framework supports hybrid static & server rendering.',
    link: 'https://nextjs.org/',
    category: 'frontend',
  },
];
