import { Card, Flex, HStack, Icon, IconButton, Kbd, Textarea } from '@chakra-ui/react';
import { useTheme } from 'next-themes';
import React from 'react';
import { useEffect, useRef, useState } from 'react';
import { LuBot, LuSend, LuChevronDown, LuPackage } from 'react-icons/lu';
import type { AiRelayOptions } from '~/api/bloefish/shared.types';
import { Button } from '~/components/ui/button';
import { Tooltip } from '~/components/ui/tooltip';
import { useChatInput } from './hooks/use-chat-input';
import { useAppDispatch } from '~/store';
import { updatePrompt } from './store';
import { StatusIndicator } from './components/atoms/StatusIndicator';
import { AiProviderPicker } from './components/molecules/AiProviderPicker';
import { SkillSetPicker } from './components/molecules/SkillSetPicker';
import { FilePicker } from './components/molecules/FilePicker';

export interface InvokedEvent {
	prompt: string;
	fileIds: string[];
	skillSetIds: string[];
	aiRelayOptions: AiRelayOptions;
}

type EnterMode = 'newline' | 'send';

interface ChatInputProps {
	disabled: boolean;
	identifier: string;
	onInvoke: (event: InvokedEvent) => void;
}

export const ChatInput: React.FC<ChatInputProps> = ({
	disabled,
	identifier,
	onInvoke,
}) => {
	const dispatch = useAppDispatch();
	const {
		prompt,
		ready,
		destinationModel,
		skillSetIds,
		fileIds,
	} = useChatInput(identifier);

	const theme = useTheme();
	const [focused, setFocused] = useState(false);
	const inputRef = useRef<HTMLTextAreaElement>(null);

	const [enterMode, setEnterMode] = useState<EnterMode>('send');

	function invoke() {
		if (disabled || !ready || !destinationModel) return;

		onInvoke({
			prompt,
			skillSetIds,
			fileIds,
			aiRelayOptions: {
				modelId: destinationModel.modelId,
				providerId: destinationModel.providerId,
			},
		});
	}

	useEffect(() => {
		if (prompt.includes('\n')) {
			setEnterMode('newline');
		} else {
			setEnterMode('send')
		}
	}, [prompt]);

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
			w={'100%'}
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
						value={prompt}
						onChange={e => dispatch(updatePrompt({
							identifier,
							prompt: e.target.value,
						}))}
						ref={inputRef}
						onFocus={() => setFocused(true)}
						onBlur={() => setFocused(false)}
						onKeyDown={e => {
							if (e.key === 'Enter' && e.metaKey && enterMode === 'newline') {
								e.preventDefault();
								invoke();
							}

							if (e.key === 'Enter' && !e.shiftKey && enterMode === 'send') {
								e.preventDefault();
								invoke();
							}
						}}
					/>
				</HStack>
				<Flex justify={'space-between'} gap={4}>
					<Flex gap={2}>
						<SkillSetPicker identifier={identifier} disabled={disabled} />
						<FilePicker identifier={identifier} disabled={disabled} />
						<Tooltip content={'Select an asset to include in this message'}>
							<Button
								size={'2xs'}
								disabled
								variant={'outline'}
							>
								<LuPackage />
								<LuChevronDown />
							</Button>
						</Tooltip>
					</Flex>
					<Flex justify={'flex-end'} align={'center'} gap={2}>
						<StatusIndicator identifier={identifier} />
						<AiProviderPicker identifier={identifier} disabled={disabled} />

						<Tooltip content={(
							<React.Fragment>
								{'Pressing Enter will send the message. However to make composing messages easier, if '}
								{'your message spans multiple lines, then you must press '}
								<Kbd>{'⌘ + ⏎'}</Kbd>
								{' to send the message.'}
							</React.Fragment>
						)}>
							<IconButton
								aria-label={'Send message'}
								disabled={disabled || !ready}
								variant={'ghost'}
								size={'2xs'}
								onClick={invoke}
							>
								<LuSend />
							</IconButton>
						</Tooltip>
					</Flex>
				</Flex>
			</Card.Body>
		</Card.Root>
	);
};
