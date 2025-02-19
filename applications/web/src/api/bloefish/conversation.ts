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
} from './conversation.types';
import { createBaseQueryWithSnake } from './base';
import { deleteConversations, injectConversations } from '~/features/conversations/store';

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

						interactions: data.interactions.map((interaction) => ({
							conversationId: data.id,
							id: interaction.id,
							messageContent: interaction.messageContent,
							aiRelayOptions: interaction.aiRelayOptions,
							owner: interaction.owner,
							streamChannelId: interaction.streamChannelId,
							createdAt: interaction.createdAt,
							updatedAt: interaction.updatedAt,
							deletedAt: interaction.deletedAt,
							completedAt: interaction.completedAt,
						})),

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

						interactions: conversation.interactions.map((interaction) => ({
							conversationId: conversation.id,
							id: interaction.id,
							messageContent: interaction.messageContent,
							aiRelayOptions: interaction.aiRelayOptions,
							streamChannelId: `${conversation.id}/${interaction.id}`, // TODO(afr): Fetch from backend
							owner: interaction.owner,

							createdAt: interaction.createdAt,
							updatedAt: interaction.updatedAt,
							completedAt: interaction.completedAt,
							deletedAt: interaction.deletedAt,
						})),

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
	}),
});
