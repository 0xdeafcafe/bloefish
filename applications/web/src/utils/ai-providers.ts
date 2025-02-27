import type { AiProvider } from '~/api/bloefish/ai-relay.types';
import type { AiRelayOptions } from '~/api/bloefish/shared.types';

export function friendlyAiRelayOptions(aiRelayOptions: AiRelayOptions, providers: AiProvider[] | undefined) {
	const provider = providers?.find((p) => p.id === aiRelayOptions.providerId);
	const model = provider?.models.find((m) => m.id === aiRelayOptions.modelId);

	if (!provider || !model) return `${aiRelayOptions.providerId} (${aiRelayOptions.modelId})`;

	return `${provider.name} (${model.name})`;
}
