import type { Actor, AiRelayOptions } from '~/api/bloefish/shared.types';

export interface ConversationPlugin {
	conversationId: string;
}

export interface InteractionPlugin {
	interactionId: string;
}

export interface CreateConversationPayload {
	conversationId: string;
	owner: Actor;
	aiRelayOptions: AiRelayOptions;

	streamChannelId: string;
	title: string | null;

	createdAt: string;
	updatedAt: string;
}

export interface AddInteractionPayload extends ConversationPlugin, InteractionPlugin {
	streamChannelId: string;
	owner: Actor;
	messageContent: string;
	aiRelayOptions: AiRelayOptions;
	markedAsExcludedAt: string | null;

	createdAt: string;
	updatedAt: string;
	completedAt: string | null;
}

export interface AddActiveInteractionPayload extends ConversationPlugin, InteractionPlugin {
	messageContent: string;
	streamChannelId: string;
	aiRelayOptions: AiRelayOptions; // TODO(afr): this should come from the backend
	markedAsExcludedAt: string | null;

	createdAt: string;
	updatedAt: string;
	completedAt: string | null;
}

export interface AddInteractionFragmentPayload extends ConversationPlugin, InteractionPlugin {
	fragment: string;
}

export interface UpdateInteractionMessageContentPayload extends ConversationPlugin, InteractionPlugin {
	content: string;
}

export interface UpdateConversationTitlePayload extends ConversationPlugin {
	title: string;
	treatAsFragment: boolean;
}

export interface DeleteConversationsPayload {
	conversationIds: string[];
}

export interface DeleteInteractionsPayload {
	interactionIds: string[];
}

export interface UpdateInteractionExcludedStatePayload extends InteractionPlugin {
	excluded: boolean;
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

export interface Interaction {
	id: string;
	conversationId: string;
	streamChannelId: string | null;
	messageContent: string;

	markedAsExcludedAt: string | null;

	owner: Actor;
	aiRelayOptions: AiRelayOptions;

	createdAt: string;
	updatedAt: string;
	completedAt: string | null;
	deletedAt: string | null;
}
