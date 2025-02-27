import { Card, Flex, HStack, Icon, IconButton, Kbd, Status, Textarea } from '@chakra-ui/react';
import { useTheme } from 'next-themes';
import React from 'react';
import { useEffect, useRef, useState } from 'react';
import { LuBot, LuSend, LuChevronDown, LuGraduationCap, LuPackage } from 'react-icons/lu';
import { aiRelayApi } from '~/api/bloefish/ai-relay';
import type { AiModel, AiProvider } from '~/api/bloefish/ai-relay.types';
import type { AiRelayOptions } from '~/api/bloefish/shared.types';
import { skillSetApi } from '~/api/bloefish/skill-set';
import { userApi } from '~/api/bloefish/user';
import { Button } from '~/components/ui/button';
import { MenuCheckboxItem, MenuContent, MenuItemCommand, MenuItemGroup, MenuRadioItem, MenuRadioItemGroup, MenuRoot, MenuTrigger } from '~/components/ui/menu';
import { Tooltip } from '~/components/ui/tooltip';
import { useLocalStorageState } from '~/hooks/use-local-storage-state';

interface AvailableModel {
	provider: AiProvider;
	model: AiModel;
}

const maxPromptLength = 5000;

type LengthStatus = 'green' | 'red' | 'yellow';
type EnterMode = 'newline' | 'send';

interface ChatInputProps {
	disabled: boolean;
	value: string;
	onChange: (value: string) => void;
	onAiRelayOptionsChange: (model: AiRelayOptions | null) => void;
	onSkillSetIdsChange: (skillSetIds: string[]) => void;
	onInvoke: () => void;
}

export const ChatInput: React.FC<ChatInputProps> = ({
	disabled,
	value,
	onChange,
	onAiRelayOptionsChange,
	onSkillSetIdsChange,
	onInvoke,
}) => {
	const { data: currentUser } = userApi.useGetOrCreateDefaultUserQuery();
	const { data: providers } = aiRelayApi.useListSupportedQuery()
	const {
		data: availableSkillSets,
		isSuccess: hasAvailableSkillSets,
		isFetching: isFetchingSkillSets,
	} = skillSetApi.useListSkillSetsByOwnerQuery({
		owner: {
			type: 'user',
			identifier: currentUser!.user.id,
		},
	});

	const theme = useTheme();
	const [focused, setFocused] = useState(false);
	const inputRef = useRef<HTMLTextAreaElement>(null);

	const [availableModels, setAvailableModels] = useState<AvailableModel[]>();
	const [selectedModel, setSelectedModel] = useLocalStorageState<AvailableModel>('chat_input.selected_model');

	const [selectedSkillSets, setSelectedSkillSets] = useState<string[]>([]);

	const [enterMode, setEnterMode] = useState<EnterMode>('send');
	const [questionLength, setQuestionLength] = useState(0);
	const [lengthStatusColor, setLengthStatusColor] = useState<LengthStatus>('green');

	const modelSelectLoading = !availableModels || availableModels.length === 0;

	useEffect(() => {
		onAiRelayOptionsChange(selectedModel ? {
			providerId: selectedModel.provider.id,
			modelId: selectedModel.model.id,
		} : null);
	}, [selectedModel]);

	useEffect(() => {
		onSkillSetIdsChange(selectedSkillSets);
	}, [selectedSkillSets]);

	useEffect(() => {
		if (!providers) return;

		const availableModels = coerceAvailableModels(providers.providers);

		if (providers)
			setAvailableModels(availableModels);

		if (!selectedModel && availableModels.length > 0) {
			setSelectedModel(availableModels[0]);
		} else if (selectedModel && !availableModels.find(model => model.provider.id === selectedModel.provider.id && model.model.id === selectedModel.model.id)) {
			setSelectedModel(availableModels[0]);
		}
	}, [providers]);

	useEffect(() => {
		if (value.includes('\n')) {
			setEnterMode('newline');
		} else {
			setEnterMode('send')
		}

		if (value.length > maxPromptLength - 500) {
			setLengthStatusColor('red');
		} else if (value.length > maxPromptLength - 1000) {
			setLengthStatusColor('yellow');
		} else {
			setLengthStatusColor('green');
		}

		if (value.length > maxPromptLength) {
			onChange(value.slice(0, maxPromptLength));
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
						value={value}
						onChange={(e) => onChange(e.target.value)}
						ref={inputRef}
						onFocus={() => setFocused(true)}
						onBlur={() => setFocused(false)}
						onKeyDown={(e) => {
							if (e.key === 'Enter' && e.metaKey && enterMode === 'newline') {
								e.preventDefault();
								onInvoke();
							}

							if (e.key === 'Enter' && !e.shiftKey && enterMode === 'send') {
								e.preventDefault();
								onInvoke();
							}
						}}
					/>
				</HStack>
				<Flex justify={'space-between'} gap={4}>
					<Flex gap={2}>
						<MenuRoot onExitComplete={() => inputRef.current?.focus()} >
							<MenuTrigger asChild>
								{/* <Tooltip content={'Select a skill set to include in this message'}> */}
								<Button
									disabled={!hasAvailableSkillSets}
									loading={isFetchingSkillSets}
									size={'2xs'}
									variant={'outline'}
								>
									<LuGraduationCap />
									{selectedSkillSets.length > 0 && `(${selectedSkillSets.length})`}
									<LuChevronDown />
								</Button>
								{/* </Tooltip> */}
							</MenuTrigger>
							<MenuContent>
								<MenuItemGroup title={'Skill sets'}>
									{availableSkillSets?.skillSets?.map(ss => (
										<MenuCheckboxItem
											key={ss.id}
											value={ss.id}
											checked={selectedSkillSets.includes(ss.id)}
											onCheckedChange={(checked) => {
												if (checked) {
													setSelectedSkillSets([...selectedSkillSets, ss.id]);
												} else {
													setSelectedSkillSets(selectedSkillSets.filter((id) => id !== ss.id));
												}
											}}
										>
											{ss.name}
										</MenuCheckboxItem>
									))}
								</MenuItemGroup>
							</MenuContent>
						</MenuRoot>

						<Tooltip content={'Select a file to include in this message'}>
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
						<Status.Root colorPalette={lengthStatusColor} userSelect={'none'} fontSize={'xs'} color={'GrayText'}>
							{questionLength === maxPromptLength ? '😱 ' : <Status.Indicator />}
							{`${questionLength}/${maxPromptLength}`}
						</Status.Root>

						<MenuRoot onExitComplete={() => inputRef.current?.focus()} >
							<MenuTrigger asChild>
								<Button
									size={'2xs'}
									loading={modelSelectLoading}
									disabled={disabled || modelSelectLoading}
									variant={'outline'}
								>
									<LuBot />
									{' '}
									{selectedModel && `${selectedModel.provider.name} (${selectedModel.model.name})`}
									{' '}
									<LuChevronDown />
								</Button>
							</MenuTrigger>
							<MenuContent>
								<MenuRadioItemGroup
									value={selectedModel ? `${selectedModel.provider.id}:${selectedModel.model.id}` : ''}
									onValueChange={(e) => {
										const model = availableModels?.find((model) => `${model.provider.id}:${model.model.id}` === e.value);

										if (model) setSelectedModel(model);
									}}
								>
									{availableModels?.map((model, index) => (
										<MenuRadioItem
											value={`${model.provider.id}:${model.model.id}`}
											key={`${model.provider.id}:${model.model.id}`}
											onClick={() => (model)}
										>
											{`${model.provider.name} (${model.model.name}) `}

											{index < 9 && (
												<MenuItemCommand>{index + 1}</MenuItemCommand>
											)}
										</MenuRadioItem>
									))}
								</MenuRadioItemGroup>
							</MenuContent>
						</MenuRoot>

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
								disabled={disabled || modelSelectLoading}
								variant={'ghost'}
								size={'2xs'}
								onClick={() => onInvoke()}
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

function coerceAvailableModels(providers: AiProvider[]): AvailableModel[] {
	return providers.map((provider) => {
		return provider.models.map((model) => {
			return {
				provider,
				model,
			};
		});
	}).flat();
}
