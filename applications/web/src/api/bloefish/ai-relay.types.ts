export interface AiModel {
	providerId: string;
	modelId: string;
}

export interface EnrichedAiModel extends AiModel {
	providerName: string;
	modelName: string;
}

export interface ListSupportedResponse {
	models: EnrichedAiModel[];
}
