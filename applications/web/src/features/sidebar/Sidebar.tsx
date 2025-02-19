import { Avatar, Box, HStack, Separator, Stack, Text } from '@chakra-ui/react';
import { SidebarLink } from './components/molecules/SidebarLink';
import { LuFolderOpen, LuFolderRoot, LuGraduationCap, LuMessageCirclePlus, LuSettings2, LuWorkflow } from 'react-icons/lu';
import { useLocation } from 'react-router';
import { motion } from 'motion/react';
import { SearchButton } from '~/components/ui/search-button';
import { useAppDispatch } from '~/store';
import { openOmni } from '../omnibar/store';

const sidebarStates = [
	'new_conversation',
	'conversations',
	'still_sets',
	'projects',
	'workflows',
	'preferences',
] as const;

type SidebarButtonState = typeof sidebarStates[number];

export const Sidebar: React.FC = () => {
	const state = useSidebarLocationState();
	const indicatorTop = calculateIndicatorTop(state);
	const dispatch = useAppDispatch();

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

				<SearchButton
					marginX={6}
					marginBottom={2}
					onClick={() => dispatch(openOmni())}
				/>

				<Box position={'relative'}>
					<motion.div
						style={{
							position: 'absolute',
							left: 0,
							width: '6px',
							height: '23px',
							borderRadius: '0.125rem',
							borderTopLeftRadius: 0,
							borderBottomLeftRadius: 0,
							background: 'var(--chakra-colors-gray-emphasized)',
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
					<Stack paddingX={6}>
						<SidebarLink
							active={state === 'new_conversation'}
							content={'New conversation'}
							path={'/'}
							icon={<LuMessageCirclePlus />}
						/>
						<SidebarLink
							active={state === 'conversations'}
							content={'Conversations'}
							path={'/conversations'}
							icon={<LuFolderOpen />}
						/>
						<SidebarLink
							active={state === 'still_sets'}
							content={'Skill sets'}
							path={'/skill-sets'}
							icon={<LuGraduationCap />}
						/>
						<SidebarLink
							active={state === 'projects'}
							content={'Projects'}
							path={'/projects'}
							icon={<LuFolderRoot />}
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

	return ((index * 31) + (index * 10)) + 5;
}
