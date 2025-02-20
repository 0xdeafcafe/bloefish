import { HStack } from '@chakra-ui/react';
import { LuClipboardCopy, LuEye, LuEyeClosed, LuThumbsDown, LuThumbsUp } from 'react-icons/lu';
import type { Interaction } from '../../store/types';
import { InteractionActionButton } from '../atoms/InteractionActionButton';
import React from 'react';
import { conversationApi } from '~/api/bloefish/conversation';
import { DeleteInteractionDialog } from './DeleteInteractionDialog';
import { toaster } from '~/components/ui/toaster';

interface ConversationInteractionActionsProps {
	interaction: Interaction;
	showDeleteOnError?: boolean;
}

export const ConversationInteractionActions: React.FC<ConversationInteractionActionsProps> = ({
	interaction,
	showDeleteOnError,
}) => {
	const [updateExcludedState, excludedMutationState] = conversationApi.useUpdateInteractionExcludedStateMutation();

	const isBot = interaction.owner.type === 'bot';
	const hasErrors = interaction.errors?.length > 0;
	const hasMessageContent = interaction.messageContent !== '';
	const hasCompletedAt = interaction.completedAt !== null;

	const pending = isBot && (!hasErrors && !hasMessageContent && !hasCompletedAt);

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
			{!hasErrors && (
				<React.Fragment>
					<InteractionActionButton
						disabled={pending}
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
							disabled={excludedMutationState.isLoading || pending}
							loading={excludedMutationState.isLoading}
							tooltip={'Exclude message from AI context'}
							onClick={() => setExcludedState(true)}
						>
							<LuEye />
						</InteractionActionButton>
					) : (
						<InteractionActionButton
							disabled={excludedMutationState.isLoading || pending}
							loading={excludedMutationState.isLoading}
							tooltip={'Include message from AI context'}
							onClick={() => setExcludedState(false)}
						>
							<LuEyeClosed />
						</InteractionActionButton>
					)}
				</React.Fragment>
			)}

			{showDeleteOnError && hasErrors && (
				<DeleteInteractionDialog
					disabled={pending}
					interactionId={interaction.id}
					onDeleteSuccess={() => {
						toaster.create({
							type: 'error',
							title: 'Message deleted',
							description: 'The message has been deleted from the conversation',
						});
					}}
				/>
			)}
		</HStack>
	);
};
