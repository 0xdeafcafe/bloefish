export type ActorType = 'user' | 'bot';

export interface Actor {
	type: ActorType;
	identifier: string;
}

export interface AiRelayOptions {
	providerId: string;
	modelId: string;
}

export interface Cher {
	code: string;
	meta: Record<string, any>;
	reasons: Cher[];
}
