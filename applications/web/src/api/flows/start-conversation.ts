import { createAsyncThunk } from '@reduxjs/toolkit';
import { conversationApi } from '../bloefish/conversation';
import type { RootState } from '~/store';
import { userApi } from '../bloefish/user';
import type { useNavigate } from 'react-router';
import { addActiveInteraction, addInteraction, startConversation } from '~/features/conversations/store';
import type { Actor, AiRelayOptions } from '../bloefish/shared.types';

interface StartConversation {
	idempotencyKey: string;
	messageContent: string;
	aiRelayOptions: AiRelayOptions;
	skillSetIds: string[];
	navigate: ReturnType<typeof useNavigate>;
}

interface StartConversationReturned {
	conversationId: string;
	interactionId: string;
	streamChannelId: string;
}

export const startConversationChain = createAsyncThunk<
	StartConversationReturned,
	StartConversation,
	{ state: RootState }
>(
	'flow/start-conversation',
	async (params, { dispatch, rejectWithValue, getState }) => {
		const state = getState();
		const user = userApi.endpoints.getOrCreateDefaultUser.select()(state);

		if (!user.data?.user)
			return rejectWithValue('Invalid default user state');

		const owner: Actor = {
			type: 'user',
			identifier: user.data.user.id,
		};

		try {
			const conversation = await dispatch(conversationApi.endpoints.createConversation.initiate({
				idempotencyKey: params.idempotencyKey,
				owner,
				aiRelayOptions: params.aiRelayOptions,
			})).unwrap();

			dispatch(startConversation({
				conversationId: conversation.id,
				owner,
				aiRelayOptions: params.aiRelayOptions,

				streamChannelId: conversation.streamChannelId,
				title: conversation.title,

				createdAt: conversation.createdAt,
				updatedAt: conversation.updatedAt,
			}));

			const interaction = await dispatch(conversationApi.endpoints.createConversationMessage.initiate({
				idempotencyKey: params.idempotencyKey,
				conversationId: conversation.id,
				messageContent: params.messageContent,
				fileIds: [],
				owner,
				aiRelayOptions: params.aiRelayOptions,
				skillSetIds: params.skillSetIds,
				options: {
					useStreaming: true,
				},
			})).unwrap();

			dispatch(addInteraction({
				conversationId: conversation.id,
				interactionId: interaction.inputInteraction.id,
				streamChannelId: interaction.streamChannelId,

				markedAsExcludedAt: null,

				messageContent: params.messageContent,
				errors: [],

				owner,
				aiRelayOptions: params.aiRelayOptions,

				createdAt: interaction.inputInteraction.createdAt,
				updatedAt: interaction.inputInteraction.updatedAt,
				completedAt: interaction.inputInteraction.completedAt,
			}));
			dispatch(addActiveInteraction({
				conversationId: conversation.id,
				interactionId: interaction.responseInteraction.id,
				streamChannelId: interaction.streamChannelId,

				markedAsExcludedAt: null,

				messageContent: '',
				errors: [],

				aiRelayOptions: params.aiRelayOptions,

				createdAt: interaction.responseInteraction.createdAt,
				updatedAt: interaction.responseInteraction.updatedAt,
				completedAt: interaction.responseInteraction.completedAt,
			}));

			params.navigate(`/conversations/${conversation.id}`, { replace: false });

			return {
				conversationId: conversation.id,
				interactionId: interaction.inputInteraction.id,
				streamChannelId: interaction.streamChannelId,
			};
		} catch (err) {
			return rejectWithValue(err);
		}
	}
);
