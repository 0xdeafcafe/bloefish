import { useAppSelector } from '~/store';
import type { EnrichedAiModel } from '~/api/bloefish/ai-relay.types';

interface ReturnType {
	prompt: string;
	skillSetIds: string[];
	fileIds: string[];
	destinationModel: EnrichedAiModel | undefined;
	ready: boolean;
}

export function useChatInput(identifier: string): ReturnType {
	const chatInput = useAppSelector(state => state.chatInput[identifier]);
	if (!chatInput) {
		return {
			prompt: '',
			skillSetIds: [],
			fileIds: [],
			destinationModel: void 0,
			ready: false,
		};
	}

	const filesReady = Object.values(chatInput.files).filter(f => f.status === 'ready');
	const ready = Boolean(chatInput.destinationModel && chatInput.prompt && filesReady);

	return {
		prompt: chatInput.prompt,
		skillSetIds: chatInput.skillSetIds,
		destinationModel: chatInput.destinationModel,
		fileIds: Object.values(chatInput.files).filter(f => f.status === 'ready').map(f => f.fileId!),
		ready,
	};
}
