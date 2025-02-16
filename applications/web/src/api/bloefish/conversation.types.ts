import type { Actor, AiRelayOptions } from './shared.types';

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
