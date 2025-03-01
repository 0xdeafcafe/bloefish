import { createAsyncThunk } from "@reduxjs/toolkit";
import { conversationApi } from "../bloefish/conversation";
import type { RootState } from "~/store";
import { userApi } from "../bloefish/user";
import { addActiveInteraction, addInteraction } from "~/features/conversations/store";
import type { Actor, AiRelayOptions } from "../bloefish/shared.types";

interface ContinueConversation {
	conversationId: string;
	idempotencyKey: string;
	messageContent: string;
	aiRelayOptions: AiRelayOptions;
	skillSetIds: string[];
	fileIds: string[];
}

interface ContinueConversationReturned {
	conversationId: string;
	interactionId: string;
	streamChannelId: string;
}

export const continueConversationChain = createAsyncThunk<
	ContinueConversationReturned,
	ContinueConversation,
	{ state: RootState }
>(
	'flow/continue-conversation',
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
			const interaction = await dispatch(conversationApi.endpoints.createConversationMessage.initiate({
				idempotencyKey: params.idempotencyKey,
				conversationId: params.conversationId,
				messageContent: params.messageContent,
				fileIds: params.fileIds,
				skillSetIds: params.skillSetIds,
				owner,
				aiRelayOptions: params.aiRelayOptions,
				options: {
					useStreaming: true,
				},
			})).unwrap();

			dispatch(addInteraction({
				conversationId: params.conversationId,
				interactionId: interaction.inputInteraction.id,
				streamChannelId: interaction.streamChannelId,
				
				markedAsExcludedAt: null,

				messageContent: params.messageContent,
				errors: [],

				aiRelayOptions: params.aiRelayOptions,
				owner,

				createdAt: interaction.inputInteraction.createdAt,
				updatedAt: interaction.inputInteraction.updatedAt,
				completedAt: interaction.inputInteraction.completedAt,
			}));
			dispatch(addActiveInteraction({
				conversationId: params.conversationId,
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

			return {
				conversationId: params.conversationId,
				interactionId: interaction.inputInteraction.id,
				streamChannelId: interaction.streamChannelId,
			};
		} catch (err) {
			return rejectWithValue(err);
		}
	}
);
