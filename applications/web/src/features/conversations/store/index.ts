import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { AddActiveInteractionPayload, AddInteractionFragmentPayload, AddInteractionPayload, Conversation, CreateConversationPayload, UpdateInteractionMessageContentPayload } from './types';

const initialState: Record<string, Conversation | undefined> = {};

export const conversationsSlice = createSlice({
	name: 'conversations',
	initialState: initialState,
	reducers: {
		injectConversations: (state, { payload }: PayloadAction<Conversation[]>) => {
			for (const conversation of payload) {
				if (state[conversation.id]) {
					// TODO(afr): Attempt to inject interactions
					continue;
				}

				state[conversation.id] = conversation;
			}
		},
		startConversation: (state, { payload }: PayloadAction<CreateConversationPayload>) => {
			state[payload.conversationId] = {
				id: payload.conversationId,
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
				id: payload.interactionId,
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
				id: payload.interactionId,
				messageContent: payload.messageContent,
				aiRelayOptions: payload.aiRelayOptions,
				owner: {
					type: 'bot',
					identifier: 'open_ai',
				},
			});
		},
		addInteractionFragment: (state, { payload }: PayloadAction<AddInteractionFragmentPayload>) => {
			const conversation = state[payload.conversationId];
			if (!conversation) {
				return;
			}

			const interaction = conversation.interactions.find(i => i.id === payload.interactionId);
			if (!interaction) {
				return;
			}

			interaction.messageContent += payload.fragment;
		},
		updateInteractionMessageContent: (state, { payload }: PayloadAction<UpdateInteractionMessageContentPayload>) => {
			const conversation = state[payload.conversationId];
			if (!conversation) {
				return;
			}

			const interaction = conversation.interactions.find(i => i.id === payload.interactionId);
			if (!interaction) {
				return;
			}

			interaction.messageContent = payload.content;
		},
	},
});

export const {
	injectConversations,
	startConversation,
	addInteraction,
	addActiveInteraction,
	addInteractionFragment,
	updateInteractionMessageContent,
} = conversationsSlice.actions;
export const conversationsReducer = conversationsSlice.reducer;
