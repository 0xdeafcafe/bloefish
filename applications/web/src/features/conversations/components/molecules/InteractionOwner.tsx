import { Avatar, Badge, Circle, Float, Stack } from '@chakra-ui/react';
import { LuBot } from 'react-icons/lu';
import type { Interaction } from '../../store/types';
import React from 'react';
import { FormatDuration } from '~/components/atoms/FormatDuration';

interface InteractionOwnerProps {
	interaction: Interaction;
}

export const InteractionOwner: React.FC<InteractionOwnerProps> = ({
	interaction,
}) => {
	switch (interaction.owner.type) {
		case 'bot':
			return (
				<Stack gap={2}>
					<Avatar.Root colorPalette="pink" variant="subtle">
						<LuBot />
						<Float placement="bottom-end" offsetX="1" offsetY="1">
							<Circle
								bg="green.500"
								size="8px"
								outline="0.2em solid"
								outlineColor="bg"
							/>
						</Float>
					</Avatar.Root>
					<Badge
						w={'fit-content'}
						variant={'outline'}
						colorPalette={'pink'}
						size={'xs'}
					>
						{`${interaction.aiRelayOptions.providerId} (${interaction.aiRelayOptions.modelId})`}
					</Badge>
					<FormatDuration
						start={interaction.createdAt}
						fontSize={'xs'}
						color={'InfoText'}
					/>
				</Stack>
			);
		
		case 'user':
			return (
				<Avatar.Root colorPalette="blue" variant="subtle">
					<Avatar.Fallback name="Alexander Forbes-Reed" />
				</Avatar.Root>
			);
		
		default: return null;
	}
};
