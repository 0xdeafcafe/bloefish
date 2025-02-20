import { Box, Center, Container, Flex, Grid, GridItem, HStack, Spinner, Stack, ButtonGroup, IconButton, Breadcrumb, Skeleton } from '@chakra-ui/react';
import { LuMailQuestion, LuSlash, LuTrash2 } from 'react-icons/lu';
import { useAppDispatch, useAppSelector } from '~/store';
import { Link, useParams } from 'react-router';
import { NotFound } from '~/pages/NotFound';
import { ChatInput } from '../chat-input/ChatInput';
import { useMemo, useState } from 'react';
import { motion } from 'motion/react';
import { continueConversationChain } from '~/api/flows/continue-conversation';
import { Helmet } from 'react-helmet-async';
import { conversationApi } from '~/api/bloefish/conversation';
import { userApi } from '~/api/bloefish/user';
import { ConversationInteraction } from './components/organisms/ConversationInteraction';
import { Tooltip } from '~/components/ui/tooltip';
import { useIdempotencyKey } from '~/hooks/useIdempotencyKey';
import { Panel } from '~/components/atoms/Panel';
import type { AiRelayOptions } from '~/api/bloefish/shared.types';

export const Conversation: React.FC = () => {
	const { conversationId } = useParams();
	const conversation = useAppSelector(s => s.conversations[conversationId ?? 'kut']);
	const dispatch = useAppDispatch();
	const { data: userData } = userApi.useGetOrCreateDefaultUserQuery();
	const { isFetching: convoFetching, isLoading: convoLoading } = conversationApi.useListConversationsWithInteractionsQuery({
		owner: {
			type: 'user',
			identifier: userData!.user.id,
		},
	}, { skip: true });

	const [idempotencyKey, generateNewIdempotencyKey] = useIdempotencyKey();
	const [aiRelayOptions, setAiRelayOptions] = useState<AiRelayOptions | null>(null);
	const [question, setQuestion] = useState('');
	const [working, setWorking] = useState(false);

	const sortedInteractions = useMemo(() => {
		if (!conversation) {
			return [];
		}

		return Object.values(conversation.interactions).sort((a, b) => {
			return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
		});
	}, [conversation?.interactions]);

	if (!conversation) {
		if (convoFetching || convoLoading) {
			return (
				<Flex w={'full'} h={'full'} justify={'center'} align={'center'}>
					<Spinner />
				</Flex>
			);
		}

		return (
			<NotFound
				title={'Conversation not found'}
				pageTitle={'Conversation not found'}
				description={'The conversation you are looking for does not exist, or this part of the app just hasn\'t been built yet. Your guess is as good as mine.'}
				Icon={LuMailQuestion}
			/>
		);
	}

	async function askQuestion() {
		if (working || !conversation || !aiRelayOptions) return;

		setWorking(true);

		try {
			await dispatch(continueConversationChain({
				conversationId: conversation.id,
				idempotencyKey: idempotencyKey,
				messageContent: question,
				aiRelayOptions: aiRelayOptions,
			})).unwrap();

			generateNewIdempotencyKey();
			setQuestion('');
		} finally {
			setWorking(false);
		}
	}

	return (
		<Box
			position={'relative'}
			width={'full'}
			height={'full'}
		>
			<Helmet>
				<title>{`${conversation.title || 'New conversation'} | Bloefish`}</title>
			</Helmet>

			<HStack
				position={'absolute'}
				bottom={0} left={0} right={0}
				justify={'center'}
				gap={0}
				filter={'blur(80px)'}
			>
				<MotionBox
					animate={{ height: ['10px', '70px', '10px'] }}
					transition={{ duration: 10, repeat: Infinity, repeatType: 'mirror' }}
					width={'10%'} background={'pink.600'} opacity={1}
				/>
				<MotionBox
					animate={{ height: ['10px', '40px', '10px'] }}
					transition={{ duration: 9, repeat: Infinity, repeatType: 'mirror' }}
					width={'20%'} background={'orange.600'} opacity={1}
				/>
				<MotionBox
					animate={{ height: ['10px', '60px', '10px'] }}
					transition={{ duration: 6, repeat: Infinity, repeatType: 'mirror' }}
					width={'15%'} background={'purple.600'} opacity={1}
				/>
				<MotionBox
					animate={{ height: ['10px', '30px', '10px'] }}
					transition={{ duration: 15, repeat: Infinity, repeatType: 'mirror' }}
					width={'8%'} background={'yellow.600'} opacity={1}
				/>
				<MotionBox
					animate={{ height: ['10px', '55px', '10px'] }}
					transition={{ duration: 12, repeat: Infinity, repeatType: 'mirror' }}
					width={'14%'} background={'red.600'} opacity={1}
				/>
			</HStack>

			<Grid
				h="full"
				templateRows="max-content 1fr max-content"
				templateColumns="1fr"
			>
				<GridItem>
					<Panel.Header>
						<Flex gap={2} align={'center'}>
							<Breadcrumb.Root>
								<Breadcrumb.List>
									<Breadcrumb.Item>
										<Breadcrumb.Link asChild>
											<Link to={'/conversations'}>
												{'Conversations'}
											</Link>
										</Breadcrumb.Link>
									</Breadcrumb.Item>
									<Breadcrumb.Separator>
										<LuSlash />
									</Breadcrumb.Separator>
									<Breadcrumb.Item>
										<Breadcrumb.CurrentLink>
											{Boolean(conversation.title) ? conversation.title : <Skeleton variant={'shine'} w={32} height={4} />}
										</Breadcrumb.CurrentLink>
									</Breadcrumb.Item>
								</Breadcrumb.List>
							</Breadcrumb.Root>
						</Flex>
						<ButtonGroup variant={'outline'} size={'xs'}>
							<Tooltip content={'Delete conversation'}>
								<IconButton
									aria-label={'Delete conversation'}
									colorPalette={'red'}
									disabled
								>
									<LuTrash2 />
								</IconButton>
							</Tooltip>
						</ButtonGroup>
					</Panel.Header>
				</GridItem>
				<GridItem
					position={'relative'}
					overflowX={'scroll'}
				>
					<Container maxW={'5xl'} minW={'5xl'} py={10}>
						<Stack
							gap={6}
							pb={20}
						>
							{sortedInteractions.map(i => (
								<ConversationInteraction
									key={i.id}
									conversation={conversation}
									interaction={i}
								/>
							))}
						</Stack>
					</Container>
				</GridItem>
				<GridItem>
					<Center mb={6}>
						<ChatInput
							disabled={working}
							value={question}
							onChange={setQuestion}
							onAiRelayOptionsChange={setAiRelayOptions}
							onInvoke={askQuestion}
						/>
					</Center>
				</GridItem>
			</Grid>
		</Box>
	)
};

const MotionBox = motion.create(Box);
