import { Avatar, Box, Card, Center, Circle, Container, Flex, Float, Grid, GridItem, HStack, Spinner, Stack, Text } from '@chakra-ui/react';
import { LuBot, LuMailQuestion } from 'react-icons/lu';
import { useAppDispatch, useAppSelector } from '~/store';
import { styled } from 'styled-components';
import { useParams } from 'react-router';
import { NotFound } from '~/pages/NotFound';
import { ChatInput } from '../chat-input/ChatInput';
import { useState } from 'react';
import { motion } from 'motion/react';
import { continueConversationChain } from '~/api/flows/continue-conversation';
import Markdown from 'react-markdown';
import { useTheme } from 'next-themes';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { twilight, coy } from 'react-syntax-highlighter/dist/esm/styles/prism';

export const Conversation: React.FC = () => {
	const theme = useTheme();
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
						<Stack gap={6} pb={20}>
							{conversation.interactions.map((i) => (
								<Flex
									key={i.interactionId}
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
												<Text fontSize={'sm'} as={'div'}>
													<MarkdownWrapper>
														<Markdown
															disallowedElements={[]}
															components={{
																code(props) {
																	const { children, className, node, ...rest } = props
																	const match = /language-(\w+)/.exec(className || '')
																	return match ? (
																		// @ts-expect-error too lazy to fix this
																		<SyntaxHighlighter
																			{...rest}
																			PreTag={'div'}
																			children={String(children).replace(/\n$/, '')}
																			language={match[1]}
																			style={theme.resolvedTheme === 'dark' ? twilight : coy}
																		/>
																	) : (
																		<code {...rest} className={className}>
																			{children}
																		</code>
																	)
																},
																h1: (props) => <Text as="h1" fontSize="2xl" fontWeight="bold" {...props} />,
																h2: (props) => <Text as="h2" fontSize="xl" fontWeight="bold" {...props} />,
																h3: (props) => <Text as="h3" fontSize="lg" fontWeight="bold" {...props} />,
																h4: (props) => <Text as="h4" fontSize="md" fontWeight="bold" {...props} />,
																h5: (props) => <Text as="h5" fontSize="sm" fontWeight="bold" {...props} />,
																h6: (props) => <Text as="h6" fontSize="xs" fontWeight="bold" {...props} />,
															}}
														>
															{i.messageContent}
														</Markdown>
													</MarkdownWrapper>
												</Text>
											)}
										</Card.Body>
									</Card.Root>

									<Box maxW={'100px'} minW={'100px'} />
								</Flex>
							))}
						</Stack>
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

const MarkdownWrapper = styled.div`
	& > * {
		margin-bottom: 16px;
	}

	& > *:last-child {
		margin-bottom: 0;
	}
`;
