import { Card, Center, Spinner, Stack } from '@chakra-ui/react';
import type { Interaction } from '../../store/types';
import React from 'react';
import { useTheme } from 'next-themes';
import { InteractionErrors } from './InteractionErrors';
import { MarkdownRenderer } from './MarkdownRenderer';
import { ConversationInteractionActions } from '../organisms/ConversationInteractionActions';

interface InteractionContentProps {
	interaction: Interaction;
}

export const InteractionContent: React.FC<InteractionContentProps> = ({
	interaction,
}) => {
	const theme = useTheme();
	const onlyErrors = Boolean(interaction.messageContent === '' && interaction.errors?.length > 0)

	return (
		<Stack gap={2}>
			{onlyErrors ? (
				<InteractionErrors errors={interaction.errors} />
			) : (
				<Card.Root
					borderRadius={"lg"}
					zIndex={11}
					blur={'10px'}
					width={'fit-content'}
					background={theme.resolvedTheme === 'dark' ? 'rgb(17 17 17 / 40%)' : 'rgb(255 255 255 / 40%)'}
				>
					<Card.Body p={3} px={5}>
						<Stack gap={4}>
							{Boolean(interaction.messageContent) && (<MarkdownRenderer markdown={interaction.messageContent} />)}
							{!interaction.messageContent && !interaction.errors?.length && (
								<Center>
									<Spinner />
								</Center>
							)}

							<InteractionErrors errors={interaction.errors} />
						</Stack>
					</Card.Body>
				</Card.Root>
			)}

			<ConversationInteractionActions
				showDeleteOnError
				interaction={interaction}
			/>
		</Stack>
	)
};
