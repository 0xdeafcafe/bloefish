import { HStack } from '@chakra-ui/react';
import { LuClipboardCopy, LuEye, LuEyeClosed, LuThumbsDown, LuThumbsUp } from 'react-icons/lu';
import type { Conversation, Interaction } from '../../store/types';
import { InteractionActionButton } from '../atoms/InteractionActionButton';
import React from 'react';
import { conversationApi } from '~/api/bloefish/conversation';
import { DeleteInteractionDialog } from './DeleteInteractionDialog';

interface ConversationInteractionActionsProps {
	conversation: Conversation;
	interaction: Interaction;
}

export const ConversationInteractionActions: React.FC<ConversationInteractionActionsProps> = ({
	conversation,
	interaction,
}) => {
	const [updateExcludedState, excludedMutationState] = conversationApi.useUpdateInteractionExcludedStateMutation();

	async function setExcludedState(excluded: boolean) {
		if (excludedMutationState.isLoading) return;

		await updateExcludedState({
			excluded,
			interactionId: interaction.id,
		});
	}

	return (
		<HStack gap={2}>
			<InteractionActionButton
				onClick={() => navigator.clipboard.writeText(interaction.messageContent)}
				tooltip={'Copy message'}
			>
				<LuClipboardCopy />
			</InteractionActionButton>

			{interaction.owner.type === 'bot' && (
				<React.Fragment>
					<InteractionActionButton disabled tooltip={'Good message'}>
						<LuThumbsUp />
					</InteractionActionButton>
					<InteractionActionButton disabled tooltip={'Bad message'}>
						<LuThumbsDown />
					</InteractionActionButton>
				</React.Fragment>
			)}

			{interaction.markedAsExcludedAt === null ? (
				<InteractionActionButton
					disabled={excludedMutationState.isLoading}
					loading={excludedMutationState.isLoading}
					tooltip={'Exclude message from AI context'}
					onClick={() => setExcludedState(true)}
				>
					<LuEye />
				</InteractionActionButton>
			) : (
				<InteractionActionButton
					disabled={excludedMutationState.isLoading}
					loading={excludedMutationState.isLoading}
					tooltip={'Include message from AI context'}
					onClick={() => setExcludedState(false)}
				>
					<LuEyeClosed />
				</InteractionActionButton>
			)}

			<DeleteInteractionDialog interactionId={interaction.id}/>
		</HStack>
	);
};
