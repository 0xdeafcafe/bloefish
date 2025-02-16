import { Card, Flex, HStack, Icon, IconButton, Status, Textarea } from "@chakra-ui/react";
import { useEffect, useRef, useState } from "react";
import { LuBot, LuSend } from "react-icons/lu";

const maxPromptLength = 5000;

type LengthStatus = 'green' | 'red' | 'yellow';

interface ChatInputProps {
	disabled: boolean;
	value: string;
	onChange: (value: string) => void;
	onInvoke: () => void;
}

export const ChatInput: React.FC<ChatInputProps> = ({
	disabled,
	value,
	onChange,
	onInvoke,
}) => {
	const [focused, setFocused] = useState(false);
	const inputRef = useRef<HTMLTextAreaElement>(null);

	const [questionLength, setQuestionLength] = useState(0);
	const [lengthStatusColor, setLengthStatusColor] = useState<LengthStatus>('green');

	useEffect(() => {
		if (value.length > maxPromptLength) {
			onChange(value.slice(0, maxPromptLength));
		}

		if (value.length > maxPromptLength - 500) {
			setLengthStatusColor('red');
		} else if (value.length > maxPromptLength - 1000) {
			setLengthStatusColor('yellow');
		} else {
			setLengthStatusColor('green');
		}

		setQuestionLength(value.length);
	}, [value]);

	useEffect(() => {
		if (inputRef.current)
			inputRef.current.focus();
	}, []);

	return (
		<Card.Root
			variant={'outline'}
			w={'2xl'}
			borderColor={focused ? 'purple.500' : 'border'}
			blur={'10px'}
			// TODO(afr): Fix light theme
			background={'rgb(17 17 17 / 60%)'}
		>
			<Card.Body p={2}>
				<HStack alignItems={'flex-start'} pb={2} pl={1} gap={0}>
					<Icon mt={2} color={focused ? 'purple.500' : void 0}>
						<LuBot />
					</Icon>
					<Textarea
						border={'none'}
						disabled={disabled}
						autoresize
						maxH={'40'}
						placeholder={'What are you too lazy to do today?'}
						variant={'outline'}
						_focus={{ border: 'transparent', outline: 'transparent' }}
						value={value}
						onChange={(e) => onChange(e.target.value)}
						ref={inputRef}
						onFocus={() => setFocused(true)}
						onBlur={() => setFocused(false)}
						onKeyDown={(e) => {
							if (e.key === 'Enter' && e.metaKey) {
								e.preventDefault();
								onInvoke();
							}
						}}
					/>
				</HStack>
				<Flex justify={'flex-end'} align={'center'} gap={2}>
					<Status.Root colorPalette={lengthStatusColor} userSelect={'none'}>
						{questionLength === 5000 ? '😱 ' : <Status.Indicator />}
						{`${questionLength}/5000`}
					</Status.Root>
					<IconButton
						aria-label={'Send message'}
						disabled={disabled}
						variant={'ghost'}
						onClick={() => onInvoke()}
					>
						<LuSend />
					</IconButton>
				</Flex>
			</Card.Body>
		</Card.Root>
	);
};
