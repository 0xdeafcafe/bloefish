export type ActorType = 'user' | 'bot';

export interface Actor {
	type: ActorType;
	identifier: string;
}

export interface AiRelayOptions {
	providerId: string;
	modelId: string;
}
