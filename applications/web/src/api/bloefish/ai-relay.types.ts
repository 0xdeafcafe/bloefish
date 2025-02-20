export interface AiProvider {
	id: string;
	name: string;
	models: AiModel[];
}

export interface AiModel {
	id: string;
	name: string;
	description: string;
}

export interface ListSupportedResponse {
	providers: AiProvider[];
}
