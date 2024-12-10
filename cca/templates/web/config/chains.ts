import { chains } from 'chain-registry';

const chainNames = ['atomone', 'cosmoshub'];

export const chainOptions = chainNames.map(
  (chainName) => chains.find((chain) => chain.chain_name === chainName)!
);