import { Avatar, Box, HStack, Separator, Stack, Text } from '@chakra-ui/react';
import { SidebarLink } from './components/molecules/SidebarLink';
import { LuBox, LuFolderOpen, LuMessageCirclePlus, LuSettings2, LuWorkflow } from 'react-icons/lu';
import { useLocation } from 'react-router';
import { motion } from 'motion/react';

const sidebarStates = [
	'new_conversation',
	'conversations',
	'files',
	'workflows',
	'preferences',
] as const;

type SidebarButtonState = typeof sidebarStates[number];

export const Sidebar: React.FC = () => {
	const state = useSidebarLocationState();
	const indicatorTop = calculateIndicatorTop(state);

	return (
		<Box
			width='250px'
			color='white'
			paddingY={3}
		>
			<Stack>
				<HStack
					gap={3}
					paddingTop={3}
					paddingX={8}
				>
					<Avatar.Root variant={'subtle'} size={'xs'} fontSize={'md'}>
						🐡
					</Avatar.Root>
					<Text fontSize={'md'} fontWeight={'bolder'} color={'WindowText'}>
						Bloefish
					</Text>
				</HStack>
				<Separator marginY={2} />
				<Box position={'relative'}>
					<motion.div
						style={{
							position: 'absolute',
							left: 0,
							width: '5px',
							height: '23px',
							borderRadius: '0.125rem',
							borderTopLeftRadius: 0,
							borderBottomLeftRadius: 0,
							background: 'var(--chakra-colors-bg-emphasized)',
						}}
						animate={{
							top: `${indicatorTop}px`
						}}
						transition={{
							type: 'spring',
							bounce: 0.3,
							duration: 0.5
						}}
					/>
					<Stack paddingLeft={6}>
						<SidebarLink
							active={state === 'new_conversation'}
							content={'New conversation'}
							path={'/testing'}
							icon={<LuMessageCirclePlus />}
						/>
						<SidebarLink
							active={state === 'conversations'}
							content={'Conversations'}
							path={'/conversations'}
							icon={<LuFolderOpen />}
						/>
						<SidebarLink
							active={state === 'files'}
							content={'Files'}
							path={'/files'}
							icon={<LuBox />}
						/>
						<SidebarLink
							active={state === 'workflows'}
							content={'Workflows'}
							path={'/workflows'}
							icon={<LuWorkflow />}
						/>
						<SidebarLink
							active={state === 'preferences'}
							content={'Preferences'}
							path={'/preferences'}
							icon={<LuSettings2 />}
						/>
					</Stack>
				</Box>
			</Stack>
		</Box>
	);
};

function useSidebarLocationState(): SidebarButtonState {
	const loc = useLocation();

	switch (true) {
		case loc.pathname.startsWith('/conversations'):
			return 'conversations';
		case loc.pathname.startsWith('/files'):
			return 'files';
		case loc.pathname.startsWith('/workflows'):
			return 'workflows';
		case loc.pathname.startsWith('/preferences'):
			return 'preferences';

		default:
			return 'new_conversation';
	}
}

function calculateIndicatorTop(state: SidebarButtonState): number {
	const index = sidebarStates.indexOf(state);

	return ((index * 31) + (index * 8)) + 4;
}
