import type { EnrichedAiModel } from '~/api/bloefish/ai-relay.types';
import type { AiRelayOptions } from '~/api/bloefish/shared.types';

export function friendlyAiRelayOptions(aiRelayOptions: AiRelayOptions, models: EnrichedAiModel[] | undefined) {
	const model = models?.find(p => p.modelId === aiRelayOptions.modelId && p.providerId === aiRelayOptions.providerId);

	if (!model) return `${aiRelayOptions.providerId} (${aiRelayOptions.modelId})`;

	return `${model.providerName} (${model.modelName})`;
}
