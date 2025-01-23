import Image from "next/image";
import { useEffect, useState } from "react";
import { useChain, useManager } from "@cosmos-kit/react";
import { Box, Combobox, Skeleton, Stack, Text } from "@interchain-ui/react";
import { Chain, AssetList } from "@chain-registry/types";

import { useDetectBreakpoints } from "@/hooks";
import { chainStore, useChainStore } from "@/contexts";
import { chainOptions } from "@/config";
import { getSignerOptions } from "@/utils";

export const ChainDropdown = () => {
  const { selectedChain } = useChainStore();
  const { chain } = useChain(selectedChain);
  const [input, setInput] = useState<string>(chain.pretty_name??chain.chain_name);
  const { isMobile } = useDetectBreakpoints();

  const [isChainsAdded, setIsChainsAdded] = useState(false);
  const { addChains, getChainLogo } = useManager();

  // local chain config
  const chainConfig: Chain = require("../../../../chain.json");
  const assetConfig: AssetList = require("../../../../assetlist.json");

  useEffect(() => {
    if (isChainsAdded) return;

    if (chainConfig && assetConfig) {
      addChains(
        [chainConfig], 
        [assetConfig], 
        getSignerOptions(),
        {
          [chainConfig.chain_name]: {
            rpc: [chainConfig.apis.rpc[0].address],
            rest: [chainConfig.apis.rest[0].address]
          }
        },
      );
      setIsChainsAdded(true);
    }
  }, [chainConfig, assetConfig, isChainsAdded]);

  const onOpenChange = (isOpen: boolean) => {};

  const chains = isChainsAdded
    ? chainOptions.concat([chainConfig])
    : chainOptions;
    
  return (
    <Combobox
      onInputChange={(input) => {
        setInput(input);
      }}
      onOpenChange={onOpenChange}
      selectedKey={selectedChain}
      onSelectionChange={(key) => {
        const chainName = key as string | null;
        if (chainName) {
          chainStore.setSelectedChain(chainName);
        }
      }}
      inputAddonStart={
        <Box display="flex" justifyContent="center" alignItems="center" px="$4">
          {input === chain.pretty_name ? (
            <Image
              src={getChainLogo(selectedChain) ?? "/images/ignite.ico"}
              alt={chain.pretty_name}
              width={24}
              height={24}
              style={{
                borderRadius: "50%",
              }}
            />
          ) : (
            <Skeleton width="24px" height="24px" borderRadius="$full" />
          )}
        </Box>
      }
      styleProps={{
        width: isMobile ? "130px" : "260px",
      }}
    >
      {chains.map((c) => (
        <Combobox.Item key={c.chain_name} textValue={c.pretty_name}>
          <Stack
            direction="horizontal"
            space={isMobile ? "$3" : "$4"}
            attributes={{ alignItems: "center" }}
          >
            <Image
              src={getChainLogo(c.chain_name) ?? "/images/ignite.ico"}
              alt={c.pretty_name??chain.chain_name}
              width={isMobile ? 18 : 24}
              height={isMobile ? 18 : 24}
              style={{
                borderRadius: "50%",
              }}
            />
            <Text fontSize={isMobile ? "12px" : "16px"} fontWeight="500">
              {c.pretty_name}
            </Text>
          </Stack>
        </Combobox.Item>
      ))}
    </Combobox>
  );
};
