import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { AddActiveInteractionPayload, AddInteractionFragmentPayload, AddInteractionPayload, Conversation, CreateConversationPayload, DeleteConversationPayload as DeleteConversationsPayload, UpdateConversationTitlePayload, UpdateInteractionMessageContentPayload } from './types';

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
				owner: payload.owner,
				aiRelayOptions: payload.aiRelayOptions,
				
				streamChannelId: payload.streamChannelId,
				title: payload.title,

				interactions: [],
				createdAt: payload.createdAt,
				updatedAt: payload.updatedAt,
				deletedAt: null,
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
				createdAt: payload.createdAt,
				updatedAt: payload.updatedAt,
				completedAt: payload.completedAt,
				deletedAt: null,
			});
		},
		addActiveInteraction: (state, { payload }: PayloadAction<AddActiveInteractionPayload>) => {
			const conversation = state[payload.conversationId];
			if (!conversation) {
				return;
			}

			conversation.interactions.push({
				id: payload.interactionId,
				conversationId: payload.conversationId,
				owner: {
					type: 'bot',
					identifier: 'open_ai',
				},
				messageContent: payload.messageContent,
				aiRelayOptions: payload.aiRelayOptions,
				streamChannelId: conversation.streamChannelId,
				createdAt: payload.createdAt,
				updatedAt: payload.updatedAt,
				completedAt: payload.completedAt,
				deletedAt: null,
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
		deleteConversations: (state, { payload }: PayloadAction<DeleteConversationsPayload>) => {
			for (const conversationId of payload.conversationIds) {
				state[conversationId] = void 0;
			}
		},
		updateConversationTitle: (state, { payload }: PayloadAction<UpdateConversationTitlePayload>) => {
			const conversation = state[payload.conversationId];
			if (!conversation) {
				return;
			}

			if (payload.treatAsFragment) {
				if (conversation.title) {
					conversation.title += payload.title;
				} else {
					conversation.title = payload.title;
				}
			} else {
				conversation.title = payload.title;
			}
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
	deleteConversations,
	updateConversationTitle,
} = conversationsSlice.actions;
export const conversationsReducer = conversationsSlice.reducer;
