import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { AddActiveInteractionPayload, AddInteractionFragment, AddInteractionPayload, Conversation, CreateConversationPayload } from './types';

const initialState: Record<string, Conversation | undefined> = {};

export const conversationsSlice = createSlice({
	name: 'conversations',
	initialState: initialState,
	reducers: {
		injectConversations: (state, { payload }: PayloadAction<Conversation[]>) => {
			for (const conversation of payload) {
				if (state[conversation.conversationId]) {
					continue;
				}

				state[conversation.conversationId] = conversation;
			}
		},
		startConversation: (state, { payload }: PayloadAction<CreateConversationPayload>) => {
			state[payload.conversationId] = {
				conversationId: payload.conversationId,
				streamChannelId: payload.streamChannelIdPrefix,
				owner: payload.owner,
				aiRelayOptions: payload.aiRelayOptions,
				interactions: [],
			};
		},
		addInteraction: (state, { payload }: PayloadAction<AddInteractionPayload>) => {
			const conversation = state[payload.conversationId];
			if (!conversation) {
				return;
			}

			conversation.interactions.push({
				conversationId: payload.conversationId,
				interactionId: payload.interactionId,
				messageContent: payload.messageContent,
				aiRelayOptions: payload.aiRelayOptions,
				owner: payload.owner,
				streamChannelId: conversation.streamChannelId,
			});
		},
		addActiveInteraction: (state, { payload }: PayloadAction<AddActiveInteractionPayload>) => {
			const conversation = state[payload.conversationId];
			if (!conversation) {
				return;
			}

			conversation.interactions.push({
				streamChannelId: conversation.streamChannelId,
				conversationId: payload.conversationId,
				interactionId: payload.interactionId,
				messageContent: payload.messageContent,
				aiRelayOptions: payload.aiRelayOptions,
				owner: {
					type: 'bot',
					identifier: 'open_ai',
				},
			});
		},
		addInteractionFragment: (state, { payload }: PayloadAction<AddInteractionFragment>) => {
			const conversation = state[payload.conversationId];
			if (!conversation) {
				return;
			}

			const interaction = conversation.interactions.find(i => i.interactionId === payload.interactionId);
			if (!interaction) {
				return;
			}

			interaction.messageContent += payload.fragment;
		},
	},
});

export const {
	injectConversations,
	startConversation,
	addInteraction,
	addActiveInteraction,
	addInteractionFragment,
} = conversationsSlice.actions;
export const conversationsReducer = conversationsSlice.reducer;
