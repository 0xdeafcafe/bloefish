import { Avatar, Circle, Float } from '@chakra-ui/react';
import { LuBot } from 'react-icons/lu';
import type { Interaction } from '../../store/types';
import React from 'react';

interface InteractionOwnerProps {
	interaction: Interaction;
}

export const InteractionOwner: React.FC<InteractionOwnerProps> = ({
	interaction,
}) => {
	switch (interaction.owner.type) {
		case 'bot':
			return (
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
