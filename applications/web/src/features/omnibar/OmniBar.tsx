import type React from 'react';
import {
	DialogBody,
	DialogContent,
	DialogRoot,
} from '~/components/ui/dialog';
import { LuSearch, LuSearchX } from 'react-icons/lu';
import { useAppDispatch, useAppSelector } from '~/store';
import { closeOmni, openOmni, toggleOmni, updateOmniQuery } from './store';
import { useEffect, useRef, useState } from 'react';
import { EmptyState, Input, Stack, VStack } from '@chakra-ui/react';
import { InputGroup } from '~/components/ui/input-group';
import { useCommands } from './hooks/use-commands';
import Fuse, { type FuseResult } from 'fuse.js'
import type { Conversation } from '../conversations/store/types';
import { OmniGroup } from './components/molecules/OmniGroup';
import type { SearchContextItem } from './types';
import { SearchItem } from './components/molecules/SearchItem';
import { OmniButton } from './components/molecules/OmniButton';

export const OmniBar: React.FC = () => {
	const { open, query } = useAppSelector(s => s.omniBar);
	const conversations = useAppSelector(s => s.conversations);
	const dispatch = useAppDispatch();
	const inputRef = useRef<HTMLInputElement>(null);
	const [inputFocused, setInputFocused] = useState(false);

	const commands = useCommands();
	const foundConversations = Object.values(conversations).filter(Boolean) as Conversation[];

	const searchContextItems: SearchContextItem[] = [
		...commands.map(command => ({ ...command, searchContextType: 'command' })),
		...Object.values(foundConversations).flatMap(conversation => 
			conversation.interactions.map(interaction => ({ 
				...interaction,
				searchContextType: 'interaction',
			})),
		),
	];

	const fuse = new Fuse(searchContextItems, {
		keys: ['name', 'keywords', 'title', 'messageContent'],
		ignoreDiacritics: true,
		includeMatches: true,
		minMatchCharLength: 2,
		threshold: 0,
		ignoreLocation: true,
	});
	const searchResults = fuse.search(query);

	useEffect(() => {
		function onKeyDown(event: KeyboardEvent) {
			if (event.key === 'k' && event.metaKey) {
				if (open && !inputFocused) {
					inputRef.current?.focus();

					return;
				}

				dispatch(toggleOmni());
			}
		}

		window.addEventListener('keydown', onKeyDown);

		return () => {
			window.removeEventListener('keydown', onKeyDown);
		};
	}, [open, inputFocused]);

	return (
		<DialogRoot
			lazyMount
			open={open}
			closeOnEscape
			closeOnInteractOutside
			onOpenChange={(e) => dispatch(e.open ? openOmni() : closeOmni())}
		>
			<DialogContent borderRadius={'lg'}>
				<DialogBody p={0}>
					<InputGroup
						width={'full'}
						flex={'1'}
						startElement={<LuSearch />}
						borderBottomRadius={0}
						borderBottomWidth={1}
					>
						<Input
							placeholder={'Search conversations and commands'}
							borderTopRadius={'lg'}
							borderBottomRadius={0}
							value={query}
							size={'lg'}
							onChange={(e) => dispatch(updateOmniQuery(e.target.value))}
							border={'none'}
							ref={inputRef}
							onFocus={() => setInputFocused(true)}
							onBlur={() => setInputFocused(false)}
							_focus={{
								outline: 'none',
							}}
						/>
					</InputGroup>

					<Stack gap={4} p={4} pb={4}>
						{query && (
							<OmniGroup title={'Search results'}>
								{searchResults.map(r => (
									<SearchItem result={r.item} matches={r.matches} />
								))}
								{searchResults.length === 0 && (
									<EmptyState.Root size={'sm'}>
										<EmptyState.Content>
											<EmptyState.Indicator>
												<LuSearchX />
											</EmptyState.Indicator>
											<VStack textAlign="center" gap={1}>
												<EmptyState.Title>
													{'Search returned no results'}
												</EmptyState.Title>
												<EmptyState.Description>
													{'You can always try being less picky'}
												</EmptyState.Description>
											</VStack>
										</EmptyState.Content>
									</EmptyState.Root>
								)}
							</OmniGroup>
						)}

						<OmniGroup title={'Commands'}>
							{commands.map(c => (
								<OmniButton iconElement={c.icon} onClick={c.onInvoke}>
									{c.name}
								</OmniButton>
							))}
						</OmniGroup>
					</Stack>
				</DialogBody>
			</DialogContent>
		</DialogRoot>
	);
};
