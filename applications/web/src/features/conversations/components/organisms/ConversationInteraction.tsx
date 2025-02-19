import { Avatar, Box, Card, Circle, Flex, Float, Spinner } from '@chakra-ui/react';
import { LuBot } from 'react-icons/lu';
import { useTheme } from 'next-themes';
import type { Conversation, Interaction } from '../../store/types';
import { MarkdownRenderer } from '../molecules/MarkdownRenderer';

interface ConversationInteractionProps {
	conversation: Conversation;
	interaction: Interaction;
}

export const ConversationInteraction: React.FC<ConversationInteractionProps> = ({
	interaction,
}) => {
	const theme = useTheme();

	return (
		<Flex
			key={interaction.id}
			gap={5}
			direction={interaction.owner.type === 'user' ? 'row-reverse' : 'row'}
		>
			{interaction.owner.type === 'bot' ? (
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
			) : (
				<Avatar.Root colorPalette="blue" variant="subtle">
					<Avatar.Fallback name="Alexander Forbes-Reed" />
				</Avatar.Root>
			)}

			<Card.Root
				borderRadius={"lg"}
				zIndex={11}
				blur={'10px'}
				background={theme.resolvedTheme === 'dark' ? 'rgb(17 17 17 / 40%)' : 'rgb(255 255 255 / 40%)'}
			>
				<Card.Body p={3} px={5}>
					{interaction.messageContent === '' && <Spinner />}
					{interaction.messageContent !== '' && (
						<MarkdownRenderer markdown={interaction.messageContent} />
					)}
				</Card.Body>
			</Card.Root>

			<Box maxW={'100px'} minW={'100px'} />
		</Flex>
	);
};
