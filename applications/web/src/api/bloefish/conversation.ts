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
} from './conversation.types';
import { createBaseQueryWithSnake } from './base';
import { injectConversations } from '~/features/conversations/store';

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
						conversationId: data.conversationId,
						owner: data.owner,
						aiRelayOptions: data.aiRelayOptions,
						interactions: data.interactions.map((interaction) => ({
							conversationId: data.conversationId,
							interactionId: interaction.id,
							messageContent: interaction.messageContent,
							aiRelayOptions: interaction.aiRelayOptions,
							owner: interaction.owner,
							streamChannelId: `${data.conversationId}/${interaction.id}`, // TODO(afr): Fetch from backend
						})),
						streamChannelId: data.conversationId,
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
						conversationId: conversation.id,
						owner: conversation.owner,
						aiRelayOptions: conversation.aiRelayOptions,
						interactions: conversation.interactions.map((interaction) => ({
							conversationId: conversation.id,
							interactionId: interaction.id,
							messageContent: interaction.messageContent,
							aiRelayOptions: interaction.aiRelayOptions,
							owner: interaction.owner,
							streamChannelId: `${conversation.id}/${interaction.id}`, // TODO(afr): Fetch from backend
						})),
						streamChannelId: conversation.id,
					}))));
				} catch (error) {
					console.error(error);
				}
			},
		}),
	}),
});
