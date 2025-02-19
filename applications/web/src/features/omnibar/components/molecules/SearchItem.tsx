import type { SearchContextItem } from '../../types';
import { Text } from '@chakra-ui/react';
import { OmniButton, OmniLink } from './OmniButton';
import { LuMessageCircle } from 'react-icons/lu';
import type { FuseResultMatch } from 'fuse.js';
import React from 'react';

interface SearchItemProps {
	result: SearchContextItem;
	matches: ReadonlyArray<FuseResultMatch> | undefined;
}

export const SearchItem: React.FC<SearchItemProps> = ({
	result,
	matches,
}) => {
	switch (result.searchContextType) {
		case 'command': {
			return (
				<OmniButton iconElement={result.icon} onClick={result.onInvoke}>
					{renderMatchScopedTextFragment(result.name, matches)}
				</OmniButton>
			);
		}

		case 'interaction': {
			return (
				<OmniLink
					iconElement={<LuMessageCircle />}
					href={`/conversations/${encodeURIComponent(result.conversationId)}?interaction=${encodeURIComponent(result.id)}`}
				>
					{renderMatchScopedTextFragment(result.messageContent, matches)}
				</OmniLink>
			);
		}

		default: return null;
	}
};

function renderMatchScopedTextFragment(text: string, matches: ReadonlyArray<FuseResultMatch> | undefined) {
	if (!matches) return text.substring(0, 40);

	let maxRange = 0;
	let selectedMatch = matches[0];

	for (const match of matches) {
		const [start, end] = match.indices[0];

		if (end - start > maxRange) {
			maxRange = end - start;
			selectedMatch = match;
		}
	}

	const [start, end] = selectedMatch.indices[0];
	if (start === end) {
		return text.substring(0, 40);
	}

	const fragment = text.substring(start, end + 1);

	return (
		<React.Fragment>
			{text.substring(0, start)}
			<Text display={'inline'} fontWeight={'bolder'}>{fragment}</Text>
			{text.substring(end + 1)}
		</React.Fragment>
	);
}
