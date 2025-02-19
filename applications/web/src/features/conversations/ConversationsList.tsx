import { Box, Breadcrumb, Button, ButtonGroup, Center, EmptyState, Icon, Link as ChakraLink, LinkOverlay, Spinner, Table, Text, VStack } from '@chakra-ui/react';
import { useAppSelector } from '~/store';
import { Helmet } from 'react-helmet-async';
import { conversationApi } from '~/api/bloefish/conversation';
import { userApi } from '~/api/bloefish/user';
import type { Conversation } from './store/types';
import { LuFishSymbol, LuSquareArrowOutUpRight, LuSquareStack, LuX } from 'react-icons/lu';

import { Panel } from '~/components/atoms/Panel';
import { Checkbox } from '~/components/ui/checkbox';
import { useState } from 'react';
import { ActionBarContent, ActionBarRoot, ActionBarSelectionTrigger, ActionBarSeparator } from '~/components/ui/action-bar';
import { Link, useNavigate } from 'react-router';
import { DeleteConversationsDialog } from './components/organisms/DeleteConversationsDialog';
import { Skeleton } from '~/components/ui/skeleton';
import { FormatDuration } from '~/components/atoms/FormatDuration';

export const ConversationsList: React.FC = () => {
	const navigate = useNavigate();
	const { data: userData } = userApi.useGetOrCreateDefaultUserQuery();
	const conversations = useAppSelector(s => s.conversations);
	const truthyConversations = Object.values(conversations).filter(Boolean) as Conversation[];

	const [selection, setSelection] = useState<string[]>([]);
	const hasSelection = selection.length > 0;
	const indeterminate = hasSelection && selection.length < truthyConversations.length;

	const { isFetching: convoFetching, isLoading: convoLoading } = conversationApi.useListConversationsWithInteractionsQuery({
		owner: {
			type: 'user',
			identifier: userData!.user.id,
		},
	}, { skip: true });

	if (convoFetching || convoLoading) {
		return (
			<Container>
				<Spinner />
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
			<Table.ScrollArea>
				<Table.Root striped interactive stickyHeader size={'md'} width={'100%'}>
					<Table.Header>
						<Table.Row>
							<Table.ColumnHeader w={6}>
								<Checkbox
									top="1"
									aria-label="Select all conversations"
									size={'sm'}
									checked={indeterminate ? 'indeterminate' : selection.length > 0}
									onCheckedChange={(changes) => setSelection(changes.checked
										? truthyConversations.map((conversation) => conversation.id)
										: [],
									)}
								/>
							</Table.ColumnHeader>
							<Table.ColumnHeader>
								{'Title'}
							</Table.ColumnHeader>
							<Table.ColumnHeader>
								{'First message preview'}
							</Table.ColumnHeader>
							<Table.ColumnHeader>
								{'Created at'}
							</Table.ColumnHeader>
							<Table.ColumnHeader>
								{'Last updated at'}
							</Table.ColumnHeader>
							<Table.ColumnHeader>
								{'Actions'}
							</Table.ColumnHeader>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{truthyConversations.map(conversation => (
							<Table.Row key={conversation.id}>
								<Table.Cell>
									<Checkbox
										top="1"
										aria-label="Select conversation"
										checked={selection.includes(conversation.id)}
										size={'sm'}
										onCheckedChange={(changes) => setSelection((prev) =>
											changes.checked
												? [...prev, conversation.id]
												: selection.filter((name) => name !== conversation.id),
										)}
									/>
								</Table.Cell>
								<ChakraLink asChild variant={'plain'} display={'contents'}>
									<Link to={`/conversations/${conversation.id}`}>
										<Table.Cell alignContent={'center'}>
											{conversation.title ?? (
												<Skeleton variant={'shine'} w={32} height={4} />
											)}
										</Table.Cell>
										<Table.Cell>
											<Text truncate>
												{conversation.interactions.at(0)?.messageContent.substring(0, 70)}
											</Text>
										</Table.Cell>
										<Table.Cell>
											<FormatDuration pointInTime={conversation.createdAt} />
										</Table.Cell>
										<Table.Cell>
											<FormatDuration pointInTime={conversation.interactions.at(0)?.updatedAt ?? conversation.updatedAt} />
										</Table.Cell>
									</Link>
								</ChakraLink>
								<Table.Cell>
									<ButtonGroup size={'2xs'} variant={'outline'}>
										<Button asChild colorPalette="gray">
											<Link to={`/conversations/${conversation.id}`}>
												<Icon size={'xs'}>
													<LuSquareArrowOutUpRight />
												</Icon>
												{'View'}
											</Link>
										</Button>
										<DeleteConversationsDialog
											conversationIds={[conversation.id]}
											onDeleteSuccess={() => setSelection([])}
											deleteButtonSize={'2xs'}
											deleteButtonIconSize={'xs'}
											deleteButtonText={'Delete'}
										/>
									</ButtonGroup>
								</Table.Cell>
							</Table.Row>
						))}
					</Table.Body>
				</Table.Root>
			</Table.ScrollArea>

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
						onDeleteSuccess={() => setSelection([])}
						deleteButtonSize={'xs'}
						deleteButtonIconSize={'xs'}
						deleteButtonText={`Delete conversation${selection.length > 1 ? 's' : ''}`}
					/>
				</ActionBarContent>
			</ActionBarRoot>
		</Container>
	);
};

const Container: React.FC<React.PropsWithChildren> = ({ children }) => (
	<Box
		position={'relative'}
		width={'full'}
		height={'full'}
	>
		<Helmet>
			<title>{'Conversations | Bloefish'}</title>
		</Helmet>

		<Panel.Header>
			<Breadcrumb.Root>
				<Breadcrumb.List>
					<Breadcrumb.Item>
						<Breadcrumb.CurrentLink>
							{'Conversations'}
						</Breadcrumb.CurrentLink>
					</Breadcrumb.Item>
				</Breadcrumb.List>
			</Breadcrumb.Root>
		</Panel.Header>

		{children}
	</Box>
);
