import type { Actor, AiRelayOptions } from './shared.types';

export interface Interaction {
	id: string;
	owner: Actor;
	messageContent: string;
	fileIds: string[];
	aiRelayOptions: AiRelayOptions;
	streamChannelId: string;
	createdAt: string;
	updatedAt: string;
	completedAt: string | null;
	deletedAt: string | null;
}

export interface Conversation {
	id: string;
	owner: Actor;
	aiRelayOptions: AiRelayOptions;
	title: string | null;
	streamChannelId: string;
	interactions: Interaction[];
	createdAt: string;
	updatedAt: string;
	deletedAt: string | null;
}

export interface CreateConversationRequest {
	idempotencyKey: string;
	owner: Actor;
	aiRelayOptions: AiRelayOptions;
}

export interface CreateConversationResponse {
	id: string;
	owner: Actor;
	aiRelayOptions: AiRelayOptions;
	title: string | null;
	streamChannelId: string;
	createdAt: string;
	updatedAt: string;
	deletedAt: string | null;
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
	inputInteraction: Interaction;
	responseInteraction: Interaction;
	streamChannelId: string;
}

export interface GetConversationWithInteractionsRequest {
	conversationId: string;
}

export interface GetConversationWithInteractionsResponse extends Conversation {}

export interface ListConversationsWithInteractionsRequest {
	owner: Actor;
}

export interface ListConversationsWithInteractionsResponse {
	conversations: Conversation[];
}

export interface DeleteConversationsRequest {
	conversationIds: string[];
}
