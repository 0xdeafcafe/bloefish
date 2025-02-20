import { Avatar, Badge, Box, Card, Center, Circle, Flex, Float, Spinner, Stack } from '@chakra-ui/react';
import { LuBot } from 'react-icons/lu';
import { useTheme } from 'next-themes';
import type { Conversation, Interaction } from '../../store/types';
import { MarkdownRenderer } from '../molecules/MarkdownRenderer';
import React from 'react';
import { ConversationInteractionActions } from './ConversationInteractionActions';
import { FormatDuration } from '~/components/atoms/FormatDuration';
import { InteractionErrors } from '../molecules/InteractionErrors';

interface ConversationInteractionProps {
	conversation: Conversation;
	interaction: Interaction;
}

export const ConversationInteraction: React.FC<ConversationInteractionProps> = ({
	conversation,
	interaction,
}) => {
	const theme = useTheme();
	const onlyErrors = Boolean(interaction.messageContent === '' && interaction.errors.length > 0)

	return (
		<Flex
			key={interaction.id}
			gap={5}
			opacity={interaction.markedAsExcludedAt ? 0.5 : 1}
			direction={interaction.owner.type === 'user' ? 'row-reverse' : 'row'}
		>
			{interaction.owner.type === 'bot' ? (
				<Flex direction={'column'} align={'start'} gap={2}>
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
					<Badge variant={'outline'} colorPalette={'pink'} size={'xs'} mt={4}>
						{`${interaction.aiRelayOptions.providerId} (${interaction.aiRelayOptions.modelId})`}
					</Badge>
					<FormatDuration
						start={interaction.createdAt}
						fontSize={'xs'}
						color={'InfoText'}
					/>
				</Flex>
			) : (
				<Avatar.Root colorPalette="blue" variant="subtle">
					<Avatar.Fallback name="Alexander Forbes-Reed" />
				</Avatar.Root>
			)}

			<Stack maxW={'100%'}>
				{onlyErrors ? (
					<InteractionErrors errors={interaction.errors} />
				) : (
					<Card.Root
						borderRadius={"lg"}
						zIndex={11}
						blur={'10px'}
						background={theme.resolvedTheme === 'dark' ? 'rgb(17 17 17 / 40%)' : 'rgb(255 255 255 / 40%)'}
					>
						<Card.Body p={3} px={5}>
							<Stack gap={4}>
								{Boolean(interaction.messageContent) && (<MarkdownRenderer markdown={interaction.messageContent} />)}
								{!interaction.messageContent && !interaction.errors.length && (
									<Center>
										<Spinner />
									</Center>
								)}

								<InteractionErrors errors={interaction.errors} />
							</Stack>
						</Card.Body>
					</Card.Root>
				)}

				<Flex
					justify={'space-between'}
					align={'center'}
				>
					<ConversationInteractionActions
						showDeleteOnError
						conversation={conversation}
						interaction={interaction}
					/>
				</Flex>
			</Stack>

			<Box maxW={'100px'} minW={'100px'} />
		</Flex>
	);
};
