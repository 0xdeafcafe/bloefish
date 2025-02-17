import type { Actor, AiRelayOptions } from './shared.types';

export interface Interaction {
	id: string;
	owner: Actor;
	messageContent: string;
	fileIds: string[];
	aiRelayOptions: AiRelayOptions;
	createdAt: string;
	updatedAt: string | null;
	confirmedAt: string | null;
	deletedAt: string | null;
}

export interface CreateConversationRequest {
	idempotencyKey: string;
	owner: Actor;
	aiRelayOptions: AiRelayOptions;
}

export interface CreateConversationResponse {
	conversationId: string;
	streamChannelIdPrefix: string;
}

export interface CreateConversationMessageRequest {
	conversationId: string;
	idempotencyKey: string;
	messageContent: string;
	fileIds: string[];
	owner: Actor;
	aiRelayOptions: AiRelayOptions;
	options: CreateConversationMessageRequestOptions;
}

export interface CreateConversationMessageRequestOptions {
	useStreaming: boolean;
}

export interface CreateConversationMessageResponse {
	conversationId: string;
	interactionId: string;
	responseInteractionId: string;
	streamChannelId: string;
}

export interface GetConversationWithInteractionsRequest {
	conversationId: string;
}

export interface GetConversationWithInteractionsResponse {
	conversationId: string;
	owner: Actor;
	aiRelayOptions: AiRelayOptions;
	interactions: Interaction[];
	createdAt: string;
	deletedAt: string | null;
}

export interface ListConversationsWithInteractionsRequest {
	owner: Actor;
}

export interface ListConversationsWithInteractionsResponse {
	conversations: ListConversationsWithInteractionsResponseConversation[];
}

export interface ListConversationsWithInteractionsResponseConversation {
	id: string;
	owner: Actor;
	aiRelayOptions: AiRelayOptions;
	interactions: Interaction[];
	createdAt: string;
	deletedAt: string | null;
}
