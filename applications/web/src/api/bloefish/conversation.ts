import { createApi } from '@reduxjs/toolkit/query/react';
import type {
	CreateConversationMessageRequest,
	CreateConversationMessageResponse,
	CreateConversationRequest,
	CreateConversationResponse,

} from './conversation.types';
import { createBaseQueryWithSnake } from './base';

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
	}),
})
