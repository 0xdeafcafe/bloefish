import { Flex } from '@chakra-ui/react';
import type { Interaction } from '../../store/types';
import React from 'react';
import { ConversationInteractionActions } from '../organisms/ConversationInteractionActions';

interface InteractionFooterProps {
	interaction: Interaction;
}

export const InteractionFooter: React.FC<InteractionFooterProps> = ({
	interaction,
}) => (
	<Flex
		gap={2}
		maxW={'100%'}
		justify={'space-between'}
		direction={interaction.owner.type === 'bot' ? 'row' : 'row-reverse'}
	>
		<ConversationInteractionActions interaction={interaction} />
	</Flex>
);
