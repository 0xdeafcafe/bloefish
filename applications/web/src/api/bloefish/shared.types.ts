export type ActorType = 'user' | 'bot';

export interface Actor {
	type: ActorType;
	identifier: string;
}

export interface AiRelayOptions {
	providerId: string;
	modelId: string;
}

export interface BloefishError {
	code: string;
	meta: Record<string, unknown>;
	reasons: BloefishError[];
}
