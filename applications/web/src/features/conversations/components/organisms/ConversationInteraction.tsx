import { Grid, GridItem } from '@chakra-ui/react';
import type { Interaction } from '../../store/types';
import React from 'react';
import { InteractionOwner } from '../molecules/InteractionOwner';
import { InteractionContent } from '../molecules/InteractionContent';

interface ConversationInteractionProps {
	interaction: Interaction;
}

export const ConversationInteraction: React.FC<ConversationInteractionProps> = ({
	interaction,
}) => {
	const user = interaction.owner.type === 'user';

	return (
		<Grid
			key={interaction.id}
			id={interaction.id}
			gap={6}
			opacity={interaction.markedAsExcludedAt ? 0.5 : 1}
			templateColumns={user ? '100px 1fr auto' : 'auto 1fr 100px'}
		>
			{user ? (
				<React.Fragment>
					<GridItem />
					<GridItem display={'flex'} flexDirection={'row-reverse'}>
						<InteractionContent interaction={interaction} />
					</GridItem>
					<GridItem>
						<InteractionOwner interaction={interaction} />
					</GridItem>
				</React.Fragment>
			) : (
				<React.Fragment>
					<GridItem>
						<InteractionOwner interaction={interaction} />
					</GridItem>
					<GridItem>
						<InteractionContent interaction={interaction} />
					</GridItem>
					<GridItem />
				</React.Fragment>
			)}
		</Grid>
	);
};
