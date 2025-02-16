import { Box, Card, Flex, HStack, Icon, IconButton, Status, Text, Textarea } from '@chakra-ui/react';
import { useEffect, useRef, useState } from 'react';
import { LuBot, LuFish, LuSend } from 'react-icons/lu';
import { TagButton } from '~/components/ui/tag-button';
import { motion } from 'framer-motion';
import { useAppDispatch } from '~/store';
import { startConversationChain } from '~/api/flows/start-conversation';

const maxPromptLength = 5000;

const starterPrompts: string[] = [
	'What is the meaning of life?',
	'Explain the concept of a monad',
	'Write a React counter component',
	'How do I use the Go Mongo Driver?',
	'Why does Hill yap so much?',
];

export const NewConversation: React.FC = () => {
	const [timeOfDay, timeOfDayEmphasis, timeOfDayEmoji] = generateTimeOfDayText();
	const [typedPrompt, setTypedPrompt] = useState('');
	const [promptCount, setPromptCount] = useState(0);
	const [promptStatusColor, setPromptStatusColor] = useState<'green' | 'red' | 'yellow'>('green');
	const inputRef = useRef<HTMLTextAreaElement>(null);
	const [isFocused, setIsFocused] = useState(false);
	
	const dispatch = useAppDispatch();
	const [working, setWorking] = useState(false);

	async function askQuestion() {
		if (working) return;

		setWorking(true);
		dispatch(startConversationChain({
			idempotencyKey: typedPrompt,
			messageContent: typedPrompt,
		}));
	}

	useEffect(() => {
		if (typedPrompt.length > maxPromptLength) {
			setTypedPrompt(typedPrompt.slice(0, maxPromptLength));
		}

		if (typedPrompt.length > maxPromptLength - 500) {
			setPromptStatusColor('red');
		} else if (typedPrompt.length > maxPromptLength - 1000) {
			setPromptStatusColor('yellow');
		} else {
			setPromptStatusColor('green');
		}

		setPromptCount(typedPrompt.length);
	}, [typedPrompt]);

	useEffect(() => {
		if (inputRef.current) {
			inputRef.current.focus();
		}
	}, []);

	return (
		<Box
			position={'relative'}
			width={'full'}
			height={'full'}
		>
			<HStack
				position={'absolute'}
				bottom={0} left={0} right={0}
				justify={'center'}
				gap={0}
				filter={'blur(80px)'}
			>
				<MotionBox
					animate={{ height: ['40px', '100px', '40px'] }}
					transition={{ duration: 10, repeat: Infinity, repeatType: 'mirror' }}
					width={'10%'} background={'pink.600'} 
				/>
				<MotionBox
					animate={{ height: ['40px', '70px', '40px'] }}
					transition={{ duration: 9, repeat: Infinity, repeatType: 'mirror' }}
					width={'20%'} background={'orange.600'}
				/>
				<MotionBox
					animate={{ height: ['40px', '90px', '40px'] }}
					transition={{ duration: 6, repeat: Infinity, repeatType: 'mirror' }}
					width={'15%'} background={'purple.600'}
				/>
				<MotionBox
					animate={{ height: ['40px', '60px', '40px'] }}
					transition={{ duration: 15, repeat: Infinity, repeatType: 'mirror' }}
					width={'8%'} background={'yellow.600'}
				/>
				<MotionBox
					animate={{ height: ['40px', '85px', '40px'] }}
					transition={{ duration: 12, repeat: Infinity, repeatType: 'mirror' }}
					width={'14%'} background={'red.600'}
				/>
			</HStack>

			<Flex
				width={'full'}
				height={'full'}
				direction={'column'}
				align={'center'}
				maxWidth={'3xl'}
				margin={'0 auto'}
			>
				<Flex
					flexGrow={1}
					direction={'column'}
					alignItems={'center'}
					justifyContent={'center'}
					gap={'6'}
				>
					<Box
						boxShadow={'lg'}
						borderRadius={'lg'}
						padding={'4'}
					>
						<Icon
							size={'2xl'}
							color={'#ffa0f9'}
						>
							<LuFish />
						</Icon>
					</Box>

					<Text
						fontSize={'5xl'}
						fontWeight={'semibold'}
						textAlign={'center'}
						lineHeight={'shorter'}
						textShadow={'lg'}
					>
						{`${timeOfDay} `}
						<span
							style={{
								background: '-webkit-linear-gradient(left, #caaffd, #ffa0f9)',
								WebkitBackgroundClip: 'text',
								WebkitTextFillColor: 'transparent'
							}}
						>
							{timeOfDayEmphasis}
						</span>
						{` ${timeOfDayEmoji}`}
						<br />
						{'Welcome To Bloefish'}
					</Text>

					<Text
						color={'GrayText'}
						fontSize={'sm'}
						textAlign={'center'}
						maxW={'lg'}
					>
						{'You can pick a starter prompt from the list below, or you can '}
						{'ask me question below.'}
					</Text>
					
					<Flex
						gap={3}
						wrap={'wrap'}
						maxW={'xl'}
						justify={'center'}
					>
						{starterPrompts.map((prompt) => (
							<TagButton
								size={'lg'}
								key={prompt}
								borderRadius={'full'}
								onClicked={prompt => {
									setTypedPrompt(prompt);
								}}
							>
								{prompt}
							</TagButton>
						))}
					</Flex>
				</Flex>

				<Card.Root
					variant={'outline'}
					w={'2xl'}
					marginBottom={6}
					borderColor={isFocused ? 'purple.500' : 'border'}
					blur={'10px'}
					// TODO(afr): Fix light theme
					background={'rgb(17 17 17 / 60%)'}
				>
					<Card.Body p={2}>
						<HStack alignItems={'flex-start'} pb={2} pl={1} gap={0}>
							<Icon mt={2} color={isFocused ? 'purple.500' : void 0}>
								<LuBot />
							</Icon>
							<Textarea
								border={'none'}
								autoresize
								maxH={'40'}
								placeholder={'What are you too lazy to do today?'}
								variant={'outline'}
								_focus={{ border: 'transparent', outline: 'transparent' }}
								value={typedPrompt}
								onChange={(e) => setTypedPrompt(e.target.value)}
								ref={inputRef}
								onFocus={() => setIsFocused(true)}
								onBlur={() => setIsFocused(false)}
								onKeyDown={(e) => {
									if (e.key === 'Enter' && e.shiftKey) {
										e.preventDefault();
										askQuestion();
									}
								}}
							/>
						</HStack>
						<Flex justify={'flex-end'} align={'center'} gap={2}>
							<Status.Root colorPalette={promptStatusColor} userSelect={'none'}>
								{promptCount === 5000 ? '😱 ' : <Status.Indicator />}
								{`${promptCount}/5000`}
							</Status.Root>
							<IconButton
								aria-label={'Send message'}
								disabled={working}
								variant={'ghost'}
								onClick={() => askQuestion()}
							>
								<LuSend />
							</IconButton>
						</Flex>
					</Card.Body>
				</Card.Root>
			</Flex>
		</Box>
	)
};

function generateTimeOfDayText(): [string, string, string] {
	const hour = new Date().getHours();
	
	if (hour >= 1 && hour < 6) {
		return ['Still', 'Awake?', '👀'];
	} else if (hour >= 6 && hour < 12) {
		return ['Good', 'Morning' , '👋']
	} else if (hour >= 12 && hour < 17) {
		return ['Good', 'Afternoon', '🌞'];
	} else {
		return ['Good', 'Evening', '🌙'];
	}
}

const MotionBox = motion(Box);
