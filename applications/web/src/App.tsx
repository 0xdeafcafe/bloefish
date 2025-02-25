import { Grid, GridItem, HStack, Link, Text } from '@chakra-ui/react';
import { Sidebar } from './features/sidebar/Sidebar';
import { LuExternalLink } from 'react-icons/lu';
import { useStreamListener } from './api/bloefish/stream';
import { PlatformStatus } from './components/molecules/PlatformStatus';
import { Panel } from './components/atoms/Panel';

export const App: React.FC<React.PropsWithChildren> = ({ children }) => {
	useStreamListener();

	return (
		<Grid templateColumns={'max-content auto'} background={'Background'} height={'100vh'}>
			<GridItem asChild>
				<Sidebar />
			</GridItem>

			<GridItem asChild>
				<Grid
					pt={3}
					pr={3}
					templateRows={'1fr auto'}
				>
					<GridItem asChild>
						<Panel.Root>
							{children}
						</Panel.Root>
					</GridItem>
					<GridItem marginRight={1} py={3}>
						<HStack justify={'end'} gap={2}>
							<PlatformStatus />
							<Text fontSize={'xs'} color={'GrayText'}>|</Text>
							<Link href={'https://github.com/0xdeafcafe/bloefish'} target={'_blank'} rel={'noopener'} fontSize={'xs'}>
								GitHub
								<LuExternalLink />
							</Link>
							<Text fontSize={'xs'} color={'GrayText'}>|</Text>
							<Link href={'https://github.com/0xdeafcafe/bloefish/issues'} target={'_blank'} rel={'no'} fontSize={'xs'}>
								GitHub Issues
								<LuExternalLink />
							</Link>
							<Text fontSize={'xs'} color={'GrayText'}>|</Text>
							<Text fontSize={'xs'} color={'GrayText'}>Version: <b>dev</b></Text>
							<Text fontSize={'xs'} color={'GrayText'}>|</Text>
							<Text fontSize={'xs'} color={'GrayText'}>
								Made with ❤️ in Amsterdam
							</Text>
						</HStack>
					</GridItem>
				</Grid>
			</GridItem>
		</Grid>
	);
};
