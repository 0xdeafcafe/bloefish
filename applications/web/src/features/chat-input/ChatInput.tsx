import { ButtonGroup, Card, Flex, HStack, Icon, IconButton, Kbd, Status, Textarea } from '@chakra-ui/react';
import { useTheme } from 'next-themes';
import { useEffect, useRef, useState } from 'react';
import { LuBot, LuSend, LuChevronUp } from 'react-icons/lu';
import type { AiRelayOptions } from '~/api/bloefish/shared.types';
import { Button } from '~/components/ui/button';
import { MenuContent, MenuItem, MenuItemCommand, MenuRoot, MenuTrigger } from '~/components/ui/menu';

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
	const theme = useTheme();
	const [focused, setFocused] = useState(false);
	const inputRef = useRef<HTMLTextAreaElement>(null);
	const [selectedModel, setSelectedModel] = useState<AiRelayOptions>({ providerId: 'open_ai', modelId: 'gpt-4o' });

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

		function onKeyDown(event: KeyboardEvent) {
			if (event.target === inputRef.current) return;
			if (event.key !== '/') return;

			if (inputRef.current) {
				event.preventDefault();
				event.stopPropagation();
				inputRef.current.focus();
			}
		}

		window.addEventListener('keydown', onKeyDown);

		return () => {
			window.removeEventListener('keydown', onKeyDown);
		};
	}, []);

	return (
		<Card.Root
			variant={'outline'}
			w={'2xl'}
			borderColor={focused ? 'purple.500' : 'border'}
			blur={'sm'}
			background={theme.resolvedTheme === 'dark' ? 'rgb(17 17 17 / 40%)' : 'rgb(255 255 255 / 60%)'}
		>
			<Card.Body p={2}>
				<HStack alignItems={'flex-start'} pb={2} pl={1} gap={0}>
					<Icon mt={'9px'} color={focused ? 'purple.500' : void 0}>
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
					<Status.Root colorPalette={lengthStatusColor} userSelect={'none'} fontSize={'xs'} color={'GrayText'}>
						{questionLength === maxPromptLength ? '😱 ' : <Status.Indicator />}
						{`${questionLength}/${maxPromptLength}`}
					</Status.Root>

					<MenuRoot>
						<MenuTrigger asChild>
						<ButtonGroup size="2xs" attached variant="outline">
							<Button
								variant="outline"
							>
								{'OpenAI (GPT-4)'}
							</Button>
							<IconButton variant="outline">
								<LuChevronUp />
							</IconButton>
						</ButtonGroup>
						</MenuTrigger>
						<MenuContent>
							<MenuItem value="open_ai:gpt-4">
								{'OpenAI (GPT-4)'} <MenuItemCommand>{'1'}</MenuItemCommand>
							</MenuItem>
							<MenuItem value="open_ai:gpt-4o">
								{'OpenAI (GPT-4o)'} <MenuItemCommand>{'2'}</MenuItemCommand>
							</MenuItem>
							<MenuItem value="open_ai:o1">
								{'OpenAI (o1)'} <MenuItemCommand>{'3'}</MenuItemCommand>
							</MenuItem>
							<MenuItem value="open_ai:o1-mini">
								{'OpenAI (o1-mini)'} <MenuItemCommand>{'4'}</MenuItemCommand>
							</MenuItem>
							<MenuItem value="open_ai:o3-mini">
								{'OpenAI (o3-mini)'} <MenuItemCommand>{'5'}</MenuItemCommand>
							</MenuItem>
						</MenuContent>
					</MenuRoot>

					<IconButton
						aria-label={'Send message'}
						disabled={disabled}
						variant={'ghost'}
						size={'2xs'}
						onClick={() => onInvoke()}
					>
						<LuSend />
					</IconButton>
				</Flex>
			</Card.Body>
		</Card.Root>
	);
};
