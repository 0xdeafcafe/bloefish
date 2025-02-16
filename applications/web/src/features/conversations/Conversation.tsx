import { Avatar, Box, Card, Center, Circle, Container, Flex, Float, Grid, GridItem, HStack, Spinner, Stack } from '@chakra-ui/react';
import { LuBot, LuMailQuestion } from 'react-icons/lu';
import { useAppDispatch, useAppSelector } from '~/store';
import { useParams } from 'react-router';
import { NotFound } from '~/pages/NotFound';
import { ChatInput } from '../chat-input/ChatInput';
import { useState } from 'react';
import { motion } from 'motion/react';
import { continueConversationChain } from '~/api/flows/continue-conversation';
import Markdown from 'react-markdown';

export const Conversation: React.FC = () => {
	const { conversationId } = useParams();
	const conversation = useAppSelector(s => s.conversationsReducer[conversationId ?? 'kut']);
	const dispatch = useAppDispatch();

	const [question, setQuestion] = useState('');
	const [working, setWorking] = useState(false);

	if (!conversation) {
		return (
			<NotFound
				title={'Conversation not found'}
				description={'The conversation you are looking for does not exist, or this part of the app just hasn\'t been built yet. Your guess is as good as mine.'}
				Icon={LuMailQuestion}
			/>
		);
	}

	async function askQuestion() {
		if (working || !conversation) return;

		setWorking(true);

		try {
			await dispatch(continueConversationChain({
				conversationId: conversation.conversationId,
				idempotencyKey: question,
				messageContent: question,
			})).unwrap();

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
			<Grid
				h="full"
				templateRows="1fr max-content"
				templateColumns="1fr"
			>
				<GridItem overflowX={'scroll'}>
					<Container maxW={'5xl'} minW={'5xl'} py={10}>
						<Center>
							<Stack gap={6}>
								{conversation.interactions.map((i) => (
									<Flex
										flex={i.interactionId}
										gap={5}
										direction={i.owner.type === 'user' ? 'row-reverse' : 'row'}
									>
										{i.owner.type === 'bot' ? (
											<Avatar.Root colorPalette="pink" variant="subtle">
												<LuBot />
												<Float placement="bottom-end" offsetX="1" offsetY="1">
													<Circle
														bg="green.500"
														size="8px"
														outline="0.2em solid"
														outlineColor="bg"
													/>
												</Float>
											</Avatar.Root>
										) : (
											<Avatar.Root colorPalette="blue" variant="subtle">
												<Avatar.Fallback name="Alexander Forbes-Reed" />
											</Avatar.Root>
										)}

										<Card.Root borderRadius={"lg"}>
											<Card.Body p={3} px={5}>
												{i.messageContent === '' ? (
													<Spinner />
												) : (
													<Markdown>
														{i.messageContent}
													</Markdown>
												)}
											</Card.Body>
										</Card.Root>
									</Flex>
								))}
							</Stack>
						</Center>
					</Container>
				</GridItem>
				<GridItem>
					<Center
						position={'relative'}
						mb={6}
					>
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

						<ChatInput
							disabled={working}
							value={question}
							onChange={setQuestion}
							onInvoke={askQuestion}
						/>
					</Center>
				</GridItem>
			</Grid>
		</Box>
	)
};

const MotionBox = motion(Box);
