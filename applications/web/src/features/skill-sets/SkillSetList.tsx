import { Box, Card, Center, Grid, GridItem, HStack, Icon, SimpleGrid, Spinner, Stack, Text } from '@chakra-ui/react';
import { Helmet } from 'react-helmet-async';
import { userApi } from '~/api/bloefish/user';
import { useRef, useState } from 'react';

import { Panel } from '~/components/atoms/Panel';
import React from 'react';
import { skillSetApi } from '~/api/bloefish/skill-set';
import { LuMicVocal, LuSquarePlus } from 'react-icons/lu';
import { Button } from '~/components/ui/button';
import { CreateSkillSetDialog } from './components/organisms/CreateSkillSetDialog';

export const SkillSetList: React.FC = () => {
	const { data: userData } = userApi.useGetOrCreateDefaultUserQuery();
	const { isFetching: skillSetFetching, isLoading: skillSetLoading, data: skillSetData } = skillSetApi.useListSkillSetsByOwnerQuery({
		owner: {
			type: 'user',
			identifier: userData!.user.id,
		},
	});

	const [showBorder, setShowBorder] = useState(false);
	const scrollContainerRef = useRef<HTMLDivElement>(null);

	const skillSets = skillSetData?.skillSets ?? [];
	const sortedSkillSets = [...skillSets].sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());

	const handleScroll = () => {
		const scrollTop = scrollContainerRef.current?.scrollTop ?? 0;
		setShowBorder(scrollTop > 24); // 1.5rem = 24px
	};

	if (skillSetFetching || skillSetLoading) {
		return (
			<Container>
				<Center>
					<Spinner />
				</Center>
			</Container>
		);
	}

	return (
		<Container>
			<Grid
				templateRows={'auto auto 1fr'}
				maxH={'100%'}
				overflow={'hidden'}
			>
				<GridItem p={6}>
					<Stack gap={2}>
						<Card.Root
							overflow={'hidden'}
							background={'bg.emphasized'}
							backgroundSize={'100vw 500px'}
						>
							<Card.Body p={6}>
								<Box blur={'lg'}>
									<Text textStyle={'4xl'} fontWeight={'bolder'} textShadow={'2xl'}>{'Skill sets'}</Text>
									<Text textStyle={'sm'} textShadow={'md'}>
										{'Skill sets store context about subjects, skills, or instructions that '}
										{'can be used to guide AI conversations. You can optionally select one '}
										{'or more skill set to enhance the expertise of the conversation.'}
									</Text>
								</Box>
							</Card.Body>
						</Card.Root>
					</Stack>
				</GridItem>

				<GridItem 
					overflow={'auto'}
					p={6}
					ref={scrollContainerRef}
					onScroll={handleScroll}
					borderTopWidth={showBorder ? '1px' : '0'}
					borderTopColor={'border.emphasized'}
				>
					<SimpleGrid minChildWidth={'sm'} gap={6}>
						<CreateSkillSetDialog
							onSuccess={() => {
								void skillSetApi.endpoints.listSkillSetsByOwner.initiate({
									owner: {
										type: 'user',
										identifier: userData!.user.id,
									},
								});
							}}
						>
							<Button
								variant={'ghost'}
								w={'full'}
								h={'auto'}
								role={'button'}
								tabIndex={0}
								asChild
							>
								<Card.Root
									border={'1px solid'}
									borderColor={'border.emphasized'}
									flexDirection={'row'}
									overflow={'hidden'}
									w={'full'}
								>
									<Card.Body>
										<Stack justify={'center'} align={'center'} h={'full'}>
											<Icon size={'2xl'}>
												<LuSquarePlus />
											</Icon>
											<Text textStyle={'lg'} color={'CaptionText'}>
												{'Add a new skill set'}
											</Text>
										</Stack>
									</Card.Body>
								</Card.Root>
							</Button>
						</CreateSkillSetDialog>

						{sortedSkillSets.map(skillSet => (
							<Card.Root
								flexDirection={'row'}
								overflow={'hidden'}
								key={skillSet.id}
							>
								<Stack>
									<Card.Body p={4} gap={2}>
										<Card.Title>
											<Icon mr={2}>
												<LuMicVocal /> 
											</Icon>
											{skillSet.name}
										</Card.Title>
										<Card.Description truncate>
											{skillSet.description}
										</Card.Description>
										<HStack>
											<Button
												size={'2xs'}
												variant={'outline'}
												colorPalette={'gray'}
											>
												Edit
											</Button>
											<Button
												size={'2xs'}
												variant={'outline'}
												colorPalette={'red'}
											>
												Delete
											</Button>
										</HStack>
									</Card.Body>
								</Stack>
							</Card.Root>
						))}
					</SimpleGrid>
				</GridItem>
			</Grid>
		</Container>
	);
};

const Container: React.FC<React.PropsWithChildren> = ({ children }) => (
	<React.Fragment>
		<Helmet>
			<title>{'Skill sets | Bloefish'}</title>
		</Helmet>

		<Panel.Body>
			{children}
		</Panel.Body>
	</React.Fragment>
);
