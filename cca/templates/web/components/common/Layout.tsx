import Head from 'next/head';
import { Box, useColorModeValue } from '@interchain-ui/react';

import { Header } from './Header';
import { Footer } from './Footer';
import { Sidebar } from './Sidebar';
import { useDisclosure } from '@/hooks';
import styles from '@/styles/layout.module.css';

export function Layout({ children }: { children: React.ReactNode }) {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <Box
      backgroundColor={useColorModeValue('$white', '$background')}
      className={styles.layout}
    >
      <Box maxWidth="1440px" width="$full" mx="$auto" display="flex">
        <Head>
          <title>Create Cosmos App</title>
          <meta name="description" content="Generated by Ignite CCA" />
          <link rel="icon" href="/images/favicon.ico" />
        </Head>
        <Sidebar isOpen={isOpen} onClose={onClose} />
        <Box
          p="30px"
          width="$full"
          minHeight="100vh"
          display="flex"
          flexDirection="column"
        >
          <Header onOpenSidebar={onOpen} />
          <Box flex="1">{children}</Box>
          <Footer />
        </Box>
      </Box>
    </Box>
  );
}
