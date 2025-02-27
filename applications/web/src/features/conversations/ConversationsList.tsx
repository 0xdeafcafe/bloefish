import { Text, Button, Center, EmptyState, Link as ChakraLink, Icon, Spinner, VStack, Card, Grid, GridItem, Stack, CheckboxCard, Badge, Skeleton, HStack } from '@chakra-ui/react';
import { useAppSelector } from '~/store';
import { Helmet } from 'react-helmet-async';
import { conversationApi } from '~/api/bloefish/conversation';
import { userApi } from '~/api/bloefish/user';
import type { Conversation } from './store/types';
import { LuFishSymbol, LuSquareArrowOutUpRight, LuSquareStack, LuX } from 'react-icons/lu';

import { Panel } from '~/components/atoms/Panel';
import { useRef, useState } from 'react';
import { ActionBarContent, ActionBarRoot, ActionBarSelectionTrigger, ActionBarSeparator } from '~/components/ui/action-bar';
import { Link, useNavigate } from 'react-router';
import { DeleteConversationsDialog } from './components/organisms/DeleteConversationsDialog';
import { toaster } from '~/components/ui/toaster';
import React from 'react';
import { aiRelayApi } from '~/api/bloefish/ai-relay';
import { friendlyAiRelayOptions } from '~/utils/ai-providers';
import { FormatDuration } from '~/components/atoms/FormatDuration';

export const ConversationsList: React.FC = () => {
	const navigate = useNavigate();
	const conversations = useAppSelector(s => s.conversations);
	const { data: userData } = userApi.useGetOrCreateDefaultUserQuery();
	const { data: supportedModels } = aiRelayApi.useListSupportedQuery();
	const truthyConversations = Object.values(conversations).filter(Boolean) as Conversation[];
	const sortedConversations = truthyConversations.sort((a, b) => {
		return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
	});

	const [selection, setSelection] = useState<string[]>([]);
	const hasSelection = selection.length > 0;

	const { isFetching: convoFetching, isLoading: convoLoading } = conversationApi.useListConversationsWithInteractionsQuery({
		owner: {
			type: 'user',
			identifier: userData!.user.id,
		},
	}, { skip: true });

	if (convoFetching || convoLoading) {
		return (
			<Container>
				<Center h={'full'}>
					<Spinner />
				</Center>
			</Container>
		);
	}

	if (truthyConversations.length === 0) {
		return (
			<Container>
				<Center h={'full'}>
					<EmptyState.Root size={'lg'}>
						<EmptyState.Content>
							<EmptyState.Indicator>
								<LuFishSymbol />
							</EmptyState.Indicator>
							<VStack textAlign={'center'}>
								<EmptyState.Title>
									{'Something on your mind?'}
								</EmptyState.Title>
								<EmptyState.Description>
									{'Once you start asking conversations, they\'ll appear here.'}
								</EmptyState.Description>
							</VStack>
						</EmptyState.Content>
					</EmptyState.Root>
				</Center>
			</Container>
		);
	}

	return (
		<Container>
			<Stack gap={6}>
				{sortedConversations.map(conversation => (
					<CheckboxCard.Root
						key={conversation.id}
						checked={selection.includes(conversation.id)}
						onChange={e => {
							if (e.target instanceof HTMLInputElement) {
								if (e.target.checked) {
									setSelection([...selection, conversation.id]);
								} else {
									setSelection(selection.filter(id => id !== conversation.id));
								}
							}
						}}
					>
						<Grid templateColumns={'auto 1fr auto'} p={2} gap={2}>
							<GridItem>
								<CheckboxCard.HiddenInput />
								<CheckboxCard.Control>
									<CheckboxCard.Indicator />
								</CheckboxCard.Control>
							</GridItem>
							<GridItem>
								<CheckboxCard.Content>
									<CheckboxCard.Label>
										{conversation.title ? (
											<ChakraLink asChild variant={'plain'} display={'contents'}>
												<Link to={`/conversations/${conversation.id}`}>
													<Text textStyle={'md'}>{conversation.title}</Text>
												</Link>
											</ChakraLink>
										) : (
											<Skeleton width={'150px'} />
										)}
									</CheckboxCard.Label>
									<CheckboxCard.Description>
										<HStack>
											<Badge size={'xs'} colorPalette={'blue'}>
												{Object.keys(conversation.interactions).length} messages
											</Badge>
											<Badge size={'xs'} colorPalette={'purple'}>
												{friendlyAiRelayOptions(conversation.aiRelayOptions, supportedModels?.providers)}
											</Badge>
											<Badge size={'xs'} colorPalette={'gray'}>
												{'Created: '}
												<FormatDuration start={conversation.createdAt} />
											</Badge>
											<Badge size={'xs'} colorPalette={'gray'}>
												{'Last updated: '}
												<FormatDuration start={conversation.createdAt} />
											</Badge>
										</HStack>
									</CheckboxCard.Description>
								</CheckboxCard.Content>
							</GridItem>
							<GridItem>
								<CheckboxCard.Content>
									<Center h={'full'} gap={2}>
										<Button asChild colorPalette={'gray'} size={'2xs'} variant={'outline'}>
											<Link to={`/conversations/${conversation.id}`}>
												<Icon size={'xs'}>
													<LuSquareArrowOutUpRight />
												</Icon>
												{'View'}
											</Link>
										</Button>
										<DeleteConversationsDialog
											conversationIds={[conversation.id]}
											onDeleteSuccess={() => {
												toaster.create({
													type: 'error',
													title: 'Conversation deleted',
													description: 'The conversation has been deleted successfully.',
												});
											}}
											deleteButtonSize={'2xs'}
											deleteButtonIconSize={'xs'}
											deleteButtonText={'Delete'}
										/>
									</Center>
								</CheckboxCard.Content>
							</GridItem>
						</Grid>
					</CheckboxCard.Root>
				))}
			</Stack>

			<ActionBarRoot open={hasSelection}>
				<ActionBarContent>
					<ActionBarSelectionTrigger>
						{selection.length} selected
					</ActionBarSelectionTrigger>
					<ActionBarSeparator />
					<Button variant={'outline'} colorPalette={'gray'} size={'xs'} onClick={() => {
						setSelection([]);
					}}>
						<Icon size={'xs'}>
							<LuX />
						</Icon>
						Clear selection
					</Button>
					<Button variant={'outline'} colorPalette={'gray'} size={'xs'} onClick={() => {
						if (selection.length === 1) {
							navigate(`/conversations/${selection[0]}`);
							setSelection([]);
							return;
						}

						for (const convoId of selection) {
							window.open(`/conversations/${convoId}`, '_blank');
						}
					}}>
						<Icon size={'xs'}>
							<LuSquareStack />
						</Icon>
						Open {selection.length > 1 ? 'all' : ''}
					</Button>
					<DeleteConversationsDialog
						conversationIds={selection}
						onDeleteSuccess={() => {
							setSelection([]);
							toaster.create({
								type: 'error',
								title: 'Conversation deleted',
								description: 'The conversation has been deleted successfully.',
							});
						}}
						deleteButtonSize={'xs'}
						deleteButtonIconSize={'xs'}
						deleteButtonText={`Delete conversation${selection.length > 1 ? 's' : ''}`}
					/>
				</ActionBarContent>
			</ActionBarRoot>
		</Container>
	);
};

const Container: React.FC<React.PropsWithChildren> = ({ children }) => {
	const [showBorder, setShowBorder] = useState(false);
	const scrollContainerRef = useRef<HTMLDivElement>(null);

	return (
		<React.Fragment>
			<Helmet>
				<title>{'Conversations | Bloefish'}</title>
			</Helmet>

			<Panel.Body>
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
									<Text textStyle={'4xl'} fontWeight={'bolder'} textShadow={'2xl'}>{'Conversations'}</Text>
									<Text textStyle={'sm'} textShadow={'md'}>
										{'weird chats below'}
									</Text>
								</Card.Body>
							</Card.Root>
						</Stack>
					</GridItem>
	
					<GridItem
						overflow={'auto'}
						p={6}
						ref={scrollContainerRef}
						onScroll={() => {
							const scrollTop = scrollContainerRef.current?.scrollTop ?? 0;
							setShowBorder(scrollTop > 24);
						}}
						borderTopWidth={showBorder ? '1px' : '0'}
						borderTopColor={'border.emphasized'}
					>
						{children}
					</GridItem>
				</Grid>
			</Panel.Body>
		</React.Fragment>
	);
}
