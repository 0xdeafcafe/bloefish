import { HStack } from '@chakra-ui/react';
import { LuClipboardCopy, LuEye, LuEyeClosed, LuThumbsDown, LuThumbsUp } from 'react-icons/lu';
import type { Conversation, Interaction } from '../../store/types';
import { InteractionActionButton } from '../atoms/InteractionActionButton';
import React from 'react';
import { conversationApi } from '~/api/bloefish/conversation';
import { DeleteInteractionDialog } from './DeleteInteractionDialog';
import { toaster } from '~/components/ui/toaster';

interface ConversationInteractionActionsProps {
	conversation: Conversation;
	interaction: Interaction;
}

export const ConversationInteractionActions: React.FC<ConversationInteractionActionsProps> = ({
	interaction,
}) => {
	const [updateExcludedState, excludedMutationState] = conversationApi.useUpdateInteractionExcludedStateMutation();
	const active = interaction.owner.type === 'bot' && interaction.completedAt === null;

	async function setExcludedState(excluded: boolean, undoable: boolean = true) {
		if (excludedMutationState.isLoading) return;

		await updateExcludedState({
			excluded,
			interactionId: interaction.id,
		}).unwrap();

		toaster.create({
			type: excluded ? 'info' : 'success',
			title: excluded ? 'Message excluded' : 'Message included',
			description: excluded ? 'The message has been excluded from the AI context' : 'The message has been included in the AI context',
			action: undoable ? {
				label: 'Undo',
				onClick: () => setExcludedState(!excluded, false),
			} : void 0,
		});
	}

	return (
		<HStack gap={2}>
			<InteractionActionButton
				onClick={() => {
					navigator.clipboard.writeText(interaction.messageContent);
					toaster.create({
						title: 'Message copied',
						description: 'The message has been copied to your clipboard',
					});
				}}
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
					disabled={excludedMutationState.isLoading || active}
					loading={excludedMutationState.isLoading}
					tooltip={'Exclude message from AI context'}
					onClick={() => setExcludedState(true)}
				>
					<LuEye />
				</InteractionActionButton>
			) : (
				<InteractionActionButton
					disabled={excludedMutationState.isLoading || active}
					loading={excludedMutationState.isLoading}
					tooltip={'Include message from AI context'}
					onClick={() => setExcludedState(false)}
				>
					<LuEyeClosed />
				</InteractionActionButton>
			)}

			<DeleteInteractionDialog
				disabled={active}
				interactionId={interaction.id}
				onDeleteSuccess={() => {
					toaster.create({
						type: 'error',
						title: 'Message deleted',
						description: 'The message has been deleted from the conversation',
					});
				}}
			/>
		</HStack>
	);
};
