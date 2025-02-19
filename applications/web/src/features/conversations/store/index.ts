import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { AddActiveInteractionPayload, AddInteractionFragmentPayload, AddInteractionPayload, Conversation, CreateConversationPayload, DeleteConversationsPayload as DeleteConversationsPayload, DeleteInteractionsPayload, UpdateConversationTitlePayload, UpdateInteractionExcludedStatePayload, UpdateInteractionMessageContentPayload } from './types';

const initialState: Record<string, Conversation | undefined> = {};

export const conversationsSlice = createSlice({
	name: 'conversations',
	initialState: initialState,
	reducers: {
		injectConversations: (state, { payload }: PayloadAction<Conversation[]>) => {
			for (const conversation of payload) {
				if (state[conversation.id]) {
					// TODO(afr): Handle injecting missing interactions

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

				interactions: {},
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

			conversation.interactions[payload.interactionId] = {
				conversationId: payload.conversationId,
				id: payload.interactionId,
				messageContent: payload.messageContent,
				streamChannelId: conversation.streamChannelId,

				markedAsExcludedAt: payload.markedAsExcludedAt,

				owner: payload.owner,
				aiRelayOptions: payload.aiRelayOptions,

				createdAt: payload.createdAt,
				updatedAt: payload.updatedAt,
				completedAt: payload.completedAt,
				deletedAt: null,
			};
		},
		addActiveInteraction: (state, { payload }: PayloadAction<AddActiveInteractionPayload>) => {
			const conversation = state[payload.conversationId];
			if (!conversation) {
				return;
			}

			conversation.interactions[payload.interactionId] = {
				id: payload.interactionId,
				conversationId: payload.conversationId,
				messageContent: payload.messageContent,
				streamChannelId: conversation.streamChannelId,

				markedAsExcludedAt: payload.markedAsExcludedAt,

				owner: {
					type: 'bot',
					identifier: 'open_ai',
				},
				aiRelayOptions: payload.aiRelayOptions,
				
				createdAt: payload.createdAt,
				updatedAt: payload.updatedAt,
				completedAt: payload.completedAt,
				deletedAt: null,
			};
		},
		addInteractionFragment: (state, { payload }: PayloadAction<AddInteractionFragmentPayload>) => {
			const conversation = state[payload.conversationId];
			if (!conversation) {
				return;
			}

			const interaction = conversation.interactions[payload.interactionId];
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

			const interaction = conversation.interactions[payload.interactionId];
			if (!interaction) {
				return;
			}

			interaction.messageContent = payload.content;
			interaction.completedAt = new Date().toISOString();
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
		deleteInteractions: (state, { payload }: PayloadAction<DeleteInteractionsPayload>) => {
			for (const conversation of Object.values(state)) {
				if (!conversation) {
					continue;
				}

				for (const interactionId of payload.interactionIds) {
					delete conversation.interactions[interactionId];
				}
			}
		},
		updateInteractionIncludedState: (state, { payload }: PayloadAction<UpdateInteractionExcludedStatePayload>) => {
			for (const conversation of Object.values(state)) {
				if (!conversation) {
					continue;
				}

				const interaction = conversation.interactions[payload.interactionId];
				if (!interaction) {
					continue;
				}

				if (payload.excluded) {
					interaction.markedAsExcludedAt = new Date().toISOString();
				} else {
					interaction.markedAsExcludedAt = null;
				}
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
	deleteInteractions,
	updateInteractionIncludedState,
} = conversationsSlice.actions;
export const conversationsReducer = conversationsSlice.reducer;
