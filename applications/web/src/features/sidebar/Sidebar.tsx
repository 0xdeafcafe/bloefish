import { Avatar, Box, HStack, Separator, Stack, Text } from '@chakra-ui/react';
import { SidebarLink } from './components/molecules/SidebarLink';
import { LuFolderOpen, LuFolderRoot, LuGraduationCap, LuLifeBuoy, LuMessageCirclePlus, LuMonitor, LuMoon, LuPackage, LuScrollText, LuSettings2, LuSun, LuWorkflow } from 'react-icons/lu';
import { useLocation } from 'react-router';
import { motion } from 'motion/react';
import { SearchButton } from '~/components/ui/search-button';
import { useAppDispatch } from '~/store';
import { openOmni } from '../omnibar/store';
import { SegmentedControl } from '~/components/ui/segmented-control';

const sidebarTopStates = [
	'new_conversation',
	'conversations',
	'workflows',
	'still_sets',
	'projects',
	'assets',
] as const;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const sidebarBottomStates = [
	'preferences',
	'change_log',
	'help',
] as const;

type SidebarTopState = typeof sidebarTopStates[number];
type SidebarButtonState = typeof sidebarBottomStates[number];
type SidebarState = SidebarTopState | SidebarButtonState;

export const Sidebar: React.FC = () => {
	const sidebarState = useSidebarLocationState();
	const indicatorTop = calculateIndicatorTop(sidebarState);
	const dispatch = useAppDispatch();

	return (
		<Box
			width='250px'
			color='white'
			paddingY={3}
		>
			<Stack h={'full'} pb={12}>
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
							top: `${indicatorTop}px`,
							opacity: sidebarTopStates.includes(sidebarState as SidebarTopState) ? 1 : 0,
						}}
						transition={{
							type: 'spring',
							bounce: 0.3,
							duration: 0.5
						}}
					/>
					<Stack paddingX={6}>
						<SidebarLink
							active={sidebarState === 'new_conversation'}
							content={'New conversation'}
							path={'/'}
							icon={<LuMessageCirclePlus />}
						/>
						<SidebarLink
							active={sidebarState === 'conversations'}
							content={'Conversations'}
							path={'/conversations'}
							icon={<LuFolderOpen />}
						/>
						<SidebarLink
							active={sidebarState === 'workflows'}
							content={'Workflows'}
							path={'/workflows'}
							icon={<LuWorkflow />}
						/>
						<SidebarLink
							active={sidebarState === 'still_sets'}
							content={'Skill sets'}
							path={'/skill-sets'}
							icon={<LuGraduationCap />}
						/>
						<SidebarLink
							active={sidebarState === 'projects'}
							content={'Projects'}
							path={'/projects'}
							icon={<LuFolderRoot />}
						/>
						<SidebarLink
							active={sidebarState === 'assets'}
							content={'Assets'}
							path={'/assets'}
							icon={<LuPackage />}
						/>
					</Stack>
				</Box>

				<Stack marginTop={'auto'} paddingX={6}>
					<SidebarLink
						active={sidebarState === 'preferences'}
						content={'Preferences'}
						path={'/preferences'}
						icon={<LuSettings2 />}
					/>
					<SidebarLink
						active={sidebarState === 'change_log'}
						content={'Change log'}
						path={'/change-log'}
						icon={<LuScrollText />}
					/>
					<SidebarLink
						active={sidebarState === 'help'}
						content={'Help'}
						path={'/help'}
						icon={<LuLifeBuoy />}
					/>

					<SegmentedControl
						display={'none'}
						size={'sm'}
						alignSelf={'center'}
						width={'fit-content'}
						colorPalette={'purple'}
						defaultValue={'system'}
						items={[
							{
								value: 'light',
								label: (
									<LuSun />
								),
							},
							{
								value: 'dark',
								label: (
									<LuMoon />
								),
							},
							{
								value: 'system',
								label: (
									<LuMonitor />
								),
							},
						]}
					/>
				</Stack>
			</Stack>
		</Box>
	);
};

function useSidebarLocationState(): SidebarState | null {
	const loc = useLocation();

	switch (true) {
		case loc.pathname === '/':
			return 'new_conversation';
		case loc.pathname.startsWith('/conversations'):
			return 'conversations';
		case loc.pathname.startsWith('/workflows'):
			return 'workflows';
		case loc.pathname.startsWith('/skill-sets'):
			return 'still_sets';
		case loc.pathname.startsWith('/projects'):
			return 'projects';
		case loc.pathname.startsWith('/preferences'):
			return 'preferences';
		case loc.pathname.startsWith('/assets'):
			return 'assets';
		
		default: return null;
	}
}

function calculateIndicatorTop(state: SidebarState | null): number {
	if (state === null) return 0;
	if (!sidebarTopStates.includes(state as SidebarTopState)) return 0;

	const index = sidebarTopStates.indexOf(state as SidebarTopState);

	return ((index * 31) + (index * 10)) + 5;
}
