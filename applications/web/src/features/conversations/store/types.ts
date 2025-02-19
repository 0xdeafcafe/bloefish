import type { Actor, AiRelayOptions } from '~/api/bloefish/shared.types';

export interface ConversationPlugin {
	conversationId: string;
}

export interface InteractionPlugin {
	interactionId: string;
}

export interface CreateConversationPayload {
	conversationId: string;
	streamChannelIdPrefix: string;
	owner: Actor;
	aiRelayOptions: AiRelayOptions;
}

export interface AddInteractionPayload extends ConversationPlugin, InteractionPlugin {
	streamChannelId: string;
	owner: Actor;
	messageContent: string;
	aiRelayOptions: AiRelayOptions;
}

export interface AddActiveInteractionPayload extends ConversationPlugin, InteractionPlugin {
	messageContent: string;
	streamChannelId: string;
	aiRelayOptions: AiRelayOptions; // TODO(afr): this should come from the backend
}

export interface AddInteractionFragmentPayload extends ConversationPlugin, InteractionPlugin {
	fragment: string;
}

export interface UpdateInteractionMessageContentPayload extends ConversationPlugin, InteractionPlugin {
	content: string;
}

export interface Conversation {
	id: string;
	streamChannelId: string;
	owner: Actor;
	aiRelayOptions: AiRelayOptions;

	interactions: Interaction[];
}

export interface Interaction {
	id: string;
	conversationId: string;
	streamChannelId: string | null;
	owner: Actor;
	messageContent: string;
	aiRelayOptions: AiRelayOptions;
}
