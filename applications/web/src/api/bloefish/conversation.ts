import { createApi } from '@reduxjs/toolkit/query/react';
import type {
	GetConversationWithInteractionsResponse,
	CreateConversationMessageRequest,
	CreateConversationMessageResponse,
	CreateConversationRequest,
	CreateConversationResponse,
	GetConversationWithInteractionsRequest,
	ListConversationsWithInteractionsRequest,
	ListConversationsWithInteractionsResponse,
	DeleteConversationsRequest,
	DeleteInteractionsRequest,
	UpdateInteractionExcludedStateRequest,
} from './conversation.types';
import { createBaseQueryWithSnake } from './base';
import { deleteConversations, deleteInteractions, injectConversations, updateInteractionIncludedState } from '~/features/conversations/store';

export const conversationApi = createApi({
	reducerPath: 'api.bloefish.conversation',
	baseQuery: createBaseQueryWithSnake('http://svc_conversation.bloefish.local:4002/rpc/'),

	endpoints: (builder) => ({
		createConversation: builder.mutation<CreateConversationResponse, CreateConversationRequest>({
			query: (body) => ({
				url: '2025-02-12/create_conversation',
				body,
			}),
		}),

		createConversationMessage: builder.mutation<CreateConversationMessageResponse, CreateConversationMessageRequest>({
			query: (body) => ({
				url: '2025-02-12/create_conversation_message',
				body,
			}),
		}),

		getConversationWithInteractions: builder.query<GetConversationWithInteractionsResponse, GetConversationWithInteractionsRequest>({
			query: (body) => ({
				url: '2025-02-12/get_conversation_with_interactions',
				body,
			}),
			async onQueryStarted(_, { dispatch, queryFulfilled }) {
				try {
					const { data } = await queryFulfilled;

					dispatch(injectConversations([{
						id: data.id,
						owner: data.owner,
						aiRelayOptions: data.aiRelayOptions,

						title: data.title,
						streamChannelId: data.streamChannelId,

						interactions: Object.fromEntries(data.interactions.map((interaction) => [
							interaction.id,
							{
								conversationId: data.id,
								id: interaction.id,
								streamChannelId: interaction.streamChannelId,
								
								markedAsExcludedAt: interaction.markedAsExcludedAt,

								messageContent: interaction.messageContent,
								errors: interaction.errors,

								owner: interaction.owner,
								aiRelayOptions: interaction.aiRelayOptions,

								createdAt: interaction.createdAt,
								updatedAt: interaction.updatedAt,
								deletedAt: interaction.deletedAt,
								completedAt: interaction.completedAt,
							}
						])),

						createdAt: data.createdAt,
						updatedAt: data.updatedAt,
						deletedAt: data.deletedAt,
					}]));
				} catch (error) {
					console.error(error);
				}
			},
		}),

		listConversationsWithInteractions: builder.query<ListConversationsWithInteractionsResponse, ListConversationsWithInteractionsRequest>({
			query: (body) => ({
				url: '2025-02-12/list_conversations_with_interactions',
				body,
			}),
			async onQueryStarted(_, { dispatch, queryFulfilled }) {
				try {
					const { data } = await queryFulfilled;

					dispatch(injectConversations(data.conversations.map((conversation) => ({
						id: conversation.id,
						owner: conversation.owner,
						aiRelayOptions: conversation.aiRelayOptions,

						title: conversation.title,
						streamChannelId: conversation.streamChannelId,

						interactions: Object.fromEntries(conversation.interactions.map((interaction) => [
							interaction.id,
							{
								conversationId: conversation.id,
								id: interaction.id,
								streamChannelId: interaction.streamChannelId,

								markedAsExcludedAt: interaction.markedAsExcludedAt,

								messageContent: interaction.messageContent,
								errors: interaction.errors,

								aiRelayOptions: interaction.aiRelayOptions,
								owner: interaction.owner,

								createdAt: interaction.createdAt,
								updatedAt: interaction.updatedAt,
								completedAt: interaction.completedAt,
								deletedAt: interaction.deletedAt,
							}
						])),

						createdAt: conversation.createdAt,
						updatedAt: conversation.updatedAt,
						deletedAt: conversation.deletedAt,
					}))));
				} catch (error) {
					console.error(error);
				}
			},
		}),

		deleteConversations: builder.mutation<void, DeleteConversationsRequest>({
			query: (body) => ({
				url: '2025-02-12/delete_conversations',
				body,
			}),
			async onQueryStarted(req, { dispatch, queryFulfilled }) {
				try {
					await queryFulfilled;

					dispatch(deleteConversations({
						conversationIds: req.conversationIds,
					}));
				} catch (error) {
					console.error(error);
				}
			},
		}),

		deleteInteractions: builder.mutation<void, DeleteInteractionsRequest>({
			query: (body) => ({
				url: '2025-02-12/delete_interactions',
				body,
			}),
			async onQueryStarted(req, { dispatch, queryFulfilled }) {
				try {
					await queryFulfilled;

					dispatch(deleteInteractions({ interactionIds: req.interactionIds }));
				} catch (error) {
					console.error(error);
				}
			},
		}),

		updateInteractionExcludedState: builder.mutation<void, UpdateInteractionExcludedStateRequest>({
			query: (body) => ({
				url: '2025-02-12/update_interaction_excluded_state',
				body,
			}),
			async onQueryStarted(req, { dispatch, queryFulfilled }) {
				try {
					await queryFulfilled;

					dispatch(updateInteractionIncludedState({
						interactionId: req.interactionId,
						excluded: req.excluded,
					}));
				} catch (error) {
					console.error(error);
				}
			},
		}),
	}),
});
