import { Box, Button, ButtonGroup, Code, Flex, Spinner, Table, Text } from '@chakra-ui/react';
import { useAppSelector } from '~/store';
import { styled } from 'styled-components';
import { motion } from 'motion/react';
import { Helmet } from 'react-helmet-async';
import { conversationApi } from '~/api/bloefish/conversation';
import { userApi } from '~/api/bloefish/user';
import type { Conversation } from './store/types';

import { ClipboardId, ClipboardRoot } from '~/components/ui/clipboard';

export const ConversationsList: React.FC = () => {
	const { data: userData } = userApi.useGetOrCreateDefaultUserQuery();
	const conversations = useAppSelector(s => s.conversations);
	const { isFetching: convoFetching, isLoading: convoLoading } = conversationApi.useListConversationsWithInteractionsQuery({
		owner: {
			type: 'user',
			identifier: userData!.user.id,
		},
	}, { skip: true });

	// const sortedConversations = Object.values(conversations).sort((a, b) => {
	// 	const aDate = new Date(a.updatedAt);
	// 	const bDate = new Date(b.updatedAt);

	// 	return aDate < bDate ? 1 : -1;
	// });

	if (convoFetching || convoLoading) {
		return (
			<Flex w={'full'} h={'full'} justify={'center'} align={'center'}>
				<Spinner />
			</Flex>
		);
	}

	const truthyConversations = Object.values(conversations).filter(Boolean) as Conversation[];

	return (
		<Box
			position={'relative'}
			width={'full'}
			height={'full'}
		>
			<Helmet>
				<title>{'Conversations | Bloefish'}</title>
			</Helmet>

			<Table.ScrollArea>
				<Table.Root striped interactive stickyHeader>
					<Table.Header>
						<Table.Row>
							<Table.ColumnHeader>
								{'ID'}
							</Table.ColumnHeader>
							<Table.ColumnHeader>
								{'First message preview'}
							</Table.ColumnHeader>
							<Table.ColumnHeader>
								{'Author'}
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
									<Code>
										<ClipboardRoot value={conversation.id}>
											<ClipboardId color="orange.fg" textStyle="xs" />
										</ClipboardRoot>
									</Code>
								</Table.Cell>
								<Table.Cell>
									<Text truncate fontSize={'sm'}>
										{conversation.interactions.at(0)?.messageContent.substring(0, 100)}
									</Text>
								</Table.Cell>
								<Table.Cell>
									<Text truncate fontSize={'sm'}>
										{'You'}
									</Text>
								</Table.Cell>
								<Table.Cell>
									<ButtonGroup size={'xs'} variant={'outline'}>
										<Button colorPalette="gray">
											{'View'}
										</Button>
										<Button colorPalette="red">
											{'Delete'}
										</Button>
									</ButtonGroup>
								</Table.Cell>
							</Table.Row>
						))}
					</Table.Body>
				</Table.Root>
			</Table.ScrollArea>
		</Box>
	)
};

const MotionBox = motion(Box);

const MarkdownWrapper = styled.div`
	& > * {
		margin-bottom: 16px;
	}

	& > *:last-child {
		margin-bottom: 0;
	}
`;
