import { createAsyncThunk } from "@reduxjs/toolkit";
import { conversationApi } from "../bloefish/conversation";
import type { RootState } from "~/store";
import { userApi } from "../bloefish/user";

export interface StartConversationChain {
	idempotencyKey: string;
	messageContent: string;
}

export const startConversationChain = createAsyncThunk<
	{ },
	StartConversationChain,
	{ state: RootState }
>(
	'flow/start-conversation',
	async (params, { dispatch, rejectWithValue, getState }) => {
		const state = getState();
		const user = userApi.endpoints.getOrCreateDefaultUser.select()(state);

		try {
			const conversation = await dispatch(conversationApi.endpoints.createConversation.initiate({
				idempotencyKey: params.idempotencyKey,
				owner: {
					identifier: user.data!.user.id!,
					type: 'user',
				},
				aiRelayOptions: {
					providerId: 'open_ai',
					modelId: 'gpt-4',
				},
			})).unwrap();

			const interaction = await dispatch(conversationApi.endpoints.createConversationMessage.initiate({
				idempotencyKey: params.idempotencyKey,
				conversationId: conversation.conversationId,
				messageContent: params.messageContent,
				fileIds: [],
				owner: {
					identifier: user.data!.user.id!,
					type: 'user',
				},
				aiRelayOptions: {
					providerId: 'open_ai',
					modelId: 'gpt-4',
				},
				options: {
					useStreaming: true,
				},
			})).unwrap();

			// const secondResult = await dispatch(conv.endpoints.secondCall.initiate({ value: firstResult.someValue })).unwrap();

			return { };
		} catch (err) {
			return rejectWithValue(err);
		}
	}
);
