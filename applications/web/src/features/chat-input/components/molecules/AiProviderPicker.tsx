import { useEffect, useState } from 'react';
import { LuBot, LuChevronDown } from 'react-icons/lu';
import { aiRelayApi } from '~/api/bloefish/ai-relay';
import type { AiProvider } from '~/api/bloefish/ai-relay.types';
import { Button } from '~/components/ui/button';
import { MenuContent, MenuRadioItem, MenuRadioItemGroup, MenuRoot, MenuTrigger } from '~/components/ui/menu';
import { useAppDispatch } from '~/store';
import { updateDestinationModel } from '../../store';
import { useChatInput } from '../../hooks/use-chat-input';
import type { EnrichedDestinationModel } from '../../store/types';

interface AiProviderPickerProps {
	disabled: boolean;
	identifier: string;
	inputRef: React.RefObject<HTMLTextAreaElement | null>;
}

export const AiProviderPicker: React.FC<AiProviderPickerProps> = ({
	disabled,
	identifier,
	inputRef,
}) => {
	const { data: providers } = aiRelayApi.useListSupportedQuery();
	const { destinationModel } = useChatInput(identifier);
	const dispatch = useAppDispatch();

	const [availableModels, setAvailableModels] = useState<EnrichedDestinationModel[]>();
	const loading = !availableModels || availableModels.length === 0;

	useEffect(() => {
		if (!providers) return;

		const availableModels = coerceAvailableModels(providers.providers);

		if (providers)
			setAvailableModels(availableModels);

		if (!destinationModel && availableModels.length > 0) {
			dispatch(updateDestinationModel({
				identifier,
				destinationModel: availableModels[0],
			}));
		} else if (destinationModel && !availableModels.find(model => model.provider.id === destinationModel.provider.id && model.model.id === destinationModel.model.id)) {
			dispatch(updateDestinationModel({
				identifier,
				destinationModel: availableModels[0],
			}));
		}
	}, [providers]);

	return (
		<MenuRoot onExitComplete={() => inputRef.current?.focus()} >
			<MenuTrigger asChild>
				<Button
					size={'2xs'}
					loading={loading}
					disabled={disabled || loading}
					variant={'outline'}
				>
					<LuBot />
					{' '}
					{destinationModel && `${destinationModel.provider.name} (${destinationModel.model.name})`}
					{' '}
					<LuChevronDown />
				</Button>
			</MenuTrigger>
			<MenuContent>
				<MenuRadioItemGroup
					value={destinationModel ? `${destinationModel.provider.id}:${destinationModel.model.id}` : ''}
					onValueChange={(e) => {
						const model = availableModels?.find((model) => `${model.provider.id}:${model.model.id}` === e.value);

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
							value={`${model.provider.id}:${model.model.id}`}
							key={`${model.provider.id}:${model.model.id}`}
							onClick={() => (model)}
						>
							{`${model.provider.name} (${model.model.name})`}
						</MenuRadioItem>
					))}
				</MenuRadioItemGroup>
			</MenuContent>
		</MenuRoot>
	);
}

function coerceAvailableModels(providers: AiProvider[]): EnrichedDestinationModel[] {
	return providers.map((provider) => {
		return provider.models.map((model) => {
			return {
				provider,
				model,
			};
		});
	}).flat();
}
