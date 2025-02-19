import { Box, Flex, HStack, Icon, Text } from '@chakra-ui/react';
import { useState } from 'react';
import { LuFish } from 'react-icons/lu';
import { TagButton } from '~/components/ui/tag-button';
import { motion } from 'framer-motion';
import { useAppDispatch } from '~/store';
import { startConversationChain } from '~/api/flows/start-conversation';
import { useNavigate } from 'react-router';
import { ChatInput } from '../chat-input/ChatInput';
import { Helmet } from 'react-helmet-async';
import { useIdempotencyKey } from '~/hooks/useIdempotencyKey';

const starterPrompts: string[] = [
	'What is the meaning of life?',
	'Magnets, how do they work?',
	'Explain the concept of a monad',
	'Write a React counter component',
	'How do I use the Go Mongo Driver?',
	'Why does Hill yap so much?',
	'If the universe is so big, why won\'t it fight me?',
];

export const NewConversation: React.FC = () => {
	const [timeOfDay, timeOfDayEmphasis, timeOfDayEmoji] = generateTimeOfDayText();
	const [question, setQuestion] = useState('');
	const [idempotencyKey, generateNewIdempotencyKey] = useIdempotencyKey();
	const navigate = useNavigate();
	const dispatch = useAppDispatch();
	const [working, setWorking] = useState(false);

	async function askQuestion(questionOverride?: string) {
		if (working) return;

		setWorking(true);

		try {
			await dispatch(startConversationChain({
				idempotencyKey: idempotencyKey,
				messageContent: questionOverride ?? question,
				navigate,
			})).unwrap();
			generateNewIdempotencyKey();
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
				<title>{'Welcome | Bloefish'}</title>
			</Helmet>
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
						{'Pick a starter prompt from the list, or you can ask me '}
						{'question below.'}
					</Text>
					
					<Flex
						gap={3}
						wrap={'wrap'}
						maxW={'xl'}
						justify={'center'}
					>
						{starterPrompts.map((prompt) => (
							<TagButton
								size={'md'}
								key={prompt}
								borderRadius={'full'}
								onClicked={prompt => {
									setQuestion(prompt);
									askQuestion(prompt);
								}}
							>
								{prompt}
							</TagButton>
						))}
					</Flex>
				</Flex>

				<Box mb={6}>
					<ChatInput
						disabled={working}
						onChange={setQuestion}
						onInvoke={askQuestion}
						value={question}
					/>
				</Box>
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

const MotionBox = motion.create(Box);
