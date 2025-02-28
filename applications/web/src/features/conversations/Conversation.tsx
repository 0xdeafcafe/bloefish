import { Box, Center, Container, Flex, Grid, GridItem, HStack, Spinner, Stack } from '@chakra-ui/react';
import { LuMailQuestion } from 'react-icons/lu';
import { useAppDispatch, useAppSelector } from '~/store';
import { useParams } from 'react-router';
import { NotFound } from '~/pages/NotFound';
import { ChatInput } from '../chat-input/ChatInput';
import { useMemo, useState } from 'react';
import { motion } from 'motion/react';
import { continueConversationChain } from '~/api/flows/continue-conversation';
import { Helmet } from 'react-helmet-async';
import { conversationApi } from '~/api/bloefish/conversation';
import { userApi } from '~/api/bloefish/user';
import { ConversationInteraction } from './components/organisms/ConversationInteraction';
import { useIdempotencyKey } from '~/hooks/useIdempotencyKey';
import { Panel } from '~/components/atoms/Panel';
import type { AiRelayOptions } from '~/api/bloefish/shared.types';
import React from 'react';
import { HeaderCard } from '~/components/atoms/HeaderCard';

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
	const [skillSetIds, setSkillSetIds] = useState<string[]>([]);

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
				<React.Fragment>
					<Helmet>
						<title>{'Conversation loading... | Bloefish'}</title>
					</Helmet>

					<Panel.Body>
						<Flex w={'full'} h={'full'} justify={'center'} align={'center'}>
							<Spinner />
						</Flex>
					</Panel.Body>
				</React.Fragment>
			);
		}

		return (
			<React.Fragment>
				<Helmet>
					<title>{'Conversation not found | Bloefish'}</title>
				</Helmet>

				<Panel.Body>
					<NotFound
						title={'Conversation not found'}
						pageTitle={'Conversation not found'}
						description={'The conversation you are looking for does not exist, or this part of the app just hasn\'t been built yet. Your guess is as good as mine.'}
						Icon={LuMailQuestion}
					/>
				</Panel.Body>
			</React.Fragment>
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
				skillSetIds: skillSetIds,
			})).unwrap();

			generateNewIdempotencyKey();
			setQuestion('');
		} finally {
			setWorking(false);
		}
	}

	return (
		<React.Fragment>
			<Helmet>
				<title>{`${conversation.title || 'New conversation'} | Bloefish`}</title>
			</Helmet>

			<Panel.Body>
				<Grid
					flex={1}
					display={'grid'}
					templateRows={'1fr auto'}
					templateColumns={'1fr'}
					overflow={'hidden'}
				>
					<GridItem
						overflowY={'scroll'}
						overflowX={'hidden'}
						position={'relative'}
						minHeight={0}
						boxShadow={'md'}
					>
						<Container maxW={'6xl'} minW={'sm'} py={10} pb={40} w={'full'}>
							<Stack gap={6} mx={10}>
								<HeaderCard
									title={conversation.title || 'New conversation'}
									description={''}
								/>

								{sortedInteractions.map(i => <ConversationInteraction key={i.id} interaction={i} />)}
							</Stack>
						</Container>
					</GridItem>
					<GridItem position={'relative'} borderTop={'1px solid'} borderTopColor={'border.emphasized'} overflow={'hidden'}>
						<HStack
							position={'absolute'}
							bottom={0} left={0} right={0}
							justify={'center'}
							gap={0}
							filter={'blur(80px)'}
						>
							<MotionBox
								animate={{ height: ['20px', '90px', '20px'] }}
								transition={{ duration: 10, repeat: Infinity, repeatType: 'mirror' }}
								width={'10%'} background={'pink.600'} opacity={1}
							/>
							<MotionBox
								animate={{ height: ['20px', '60px', '20px'] }}
								transition={{ duration: 9, repeat: Infinity, repeatType: 'mirror' }}
								width={'20%'} background={'orange.600'} opacity={1}
							/>
							<MotionBox
								animate={{ height: ['20px', '80px', '20px'] }}
								transition={{ duration: 6, repeat: Infinity, repeatType: 'mirror' }}
								width={'15%'} background={'purple.600'} opacity={1}
							/>
							<MotionBox
								animate={{ height: ['20px', '30px', '20px'] }}
								transition={{ duration: 15, repeat: Infinity, repeatType: 'mirror' }}
								width={'8%'} background={'yellow.600'} opacity={1}
							/>
							<MotionBox
								animate={{ height: ['20px', '55px', '20px'] }}
								transition={{ duration: 12, repeat: Infinity, repeatType: 'mirror' }}
								width={'14%'} background={'red.600'} opacity={1}
							/>
						</HStack>

						<Center
							maxW={'2xl'}
							w={'full'}
							p={6}
							margin={'0 auto'}
						>
							<ChatInput
								disabled={working}
								value={question}
								onChange={setQuestion}
								onAiRelayOptionsChange={setAiRelayOptions}
								onSkillSetIdsChange={setSkillSetIds}
								onInvoke={askQuestion}
							/>
						</Center>
					</GridItem>
				</Grid>
			</Panel.Body>
		</React.Fragment>
	);
};

const MotionBox = motion.create(Box);
