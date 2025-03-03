import { useEffect, useState } from 'react';
import { LuBot, LuChevronDown } from 'react-icons/lu';
import { aiRelayApi } from '~/api/bloefish/ai-relay';
import { Button } from '~/components/ui/button';
import { MenuContent, MenuRadioItem, MenuRadioItemGroup, MenuRoot, MenuTrigger } from '~/components/ui/menu';
import { useAppDispatch } from '~/store';
import { updateDestinationModel } from '../../store';
import { useChatInput } from '../../hooks/use-chat-input';
import type { EnrichedAiModel } from '~/api/bloefish/ai-relay.types';

interface AiProviderPickerProps {
	disabled: boolean;
	identifier: string;
}

export const AiProviderPicker: React.FC<AiProviderPickerProps> = ({
	disabled,
	identifier,
}) => {
	const { data: providers } = aiRelayApi.useListSupportedQuery();
	const { destinationModel } = useChatInput(identifier);
	const dispatch = useAppDispatch();

	const [availableModels, setAvailableModels] = useState<EnrichedAiModel[]>();
	const loading = !availableModels || availableModels.length === 0;

	useEffect(() => {
		if (!providers) return;

		const availableModels = providers.models;

		if (availableModels)
			setAvailableModels(availableModels);

		if (!destinationModel && availableModels.length > 0) {
			dispatch(updateDestinationModel({
				identifier,
				destinationModel: availableModels[0],
			}));
		} else if (destinationModel && !availableModels.find(model => model.providerId === destinationModel.providerId && model.modelId === destinationModel.modelId)) {
			dispatch(updateDestinationModel({
				identifier,
				destinationModel: availableModels[0],
			}));
		}
	}, [providers]);

	return (
		<MenuRoot>
			<MenuTrigger asChild>
				<Button
					size={'2xs'}
					loading={loading}
					disabled={disabled || loading}
					variant={'outline'}
				>
					<LuBot />
					{' '}
					{destinationModel && `${destinationModel.providerName} (${destinationModel.modelName})`}
					{' '}
					<LuChevronDown />
				</Button>
			</MenuTrigger>
			<MenuContent>
				<MenuRadioItemGroup
					value={destinationModel ? `${destinationModel.providerId}:${destinationModel.modelId}` : ''}
					onValueChange={(e) => {
						const model = availableModels?.find((model) => `${model.providerId}:${model.modelId}` === e.value);

						if (model) {
							dispatch(updateDestinationModel({
								identifier,
								destinationModel: model,
							}));
						}
					}}
				>
					{availableModels?.map(model => (
						<MenuRadioItem
							value={`${model.providerId}:${model.modelId}`}
							key={`${model.providerId}:${model.modelId}`}
							onClick={() => (model)}
						>
							{`${model.providerName} (${model.modelName})`}
						</MenuRadioItem>
					))}
				</MenuRadioItemGroup>
			</MenuContent>
		</MenuRoot>
	);
}
