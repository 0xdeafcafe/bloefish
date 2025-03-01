import { Status } from '@chakra-ui/react';
import React from 'react';
import { useEffect, useState } from 'react';
import type { AiRelayOptions } from '~/api/bloefish/shared.types';
import { maxPromptLength } from '../../constants';
import { useAppDispatch } from '~/store';
import { updatePrompt } from '../../store';
import { useChatInput } from '../../hooks/use-chat-input';

export interface InvokedEvent {
	fileIds: string[];
	skillSetIds: string[];
	aiRelayOptions: AiRelayOptions;
}

type LengthStatus = 'green' | 'red' | 'yellow';

interface StatusIndicatorProps {
	identifier: string;
}

export const StatusIndicator: React.FC<StatusIndicatorProps> = ({
	identifier,
}) => {
	const { prompt } = useChatInput(identifier);
	const dispatch = useAppDispatch();

	const [questionLength, setQuestionLength] = useState(0);
	const [lengthStatusColor, setLengthStatusColor] = useState<LengthStatus>('green');

	useEffect(() => {
		if (prompt.length > maxPromptLength - 500) {
			setLengthStatusColor('red');
		} else if (prompt.length > maxPromptLength - 1000) {
			setLengthStatusColor('yellow');
		} else {
			setLengthStatusColor('green');
		}

		if (prompt.length > maxPromptLength) {
			dispatch(updatePrompt({
				identifier,
				prompt: prompt.slice(0, maxPromptLength),
			}));
		}

		setQuestionLength(prompt.length);
	}, [prompt]);

	return (
		<Status.Root colorPalette={lengthStatusColor} userSelect={'none'} fontSize={'xs'} color={'GrayText'}>
			{questionLength >= maxPromptLength ? '😱 ' : <Status.Indicator />}
			{`${questionLength}/${maxPromptLength}`}
		</Status.Root>
	);
};

