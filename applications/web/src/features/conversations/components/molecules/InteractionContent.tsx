import { Badge, Card, Center, Spinner, Stack } from '@chakra-ui/react';
import type { Interaction } from '../../store/types';
import React from 'react';
import { useTheme } from 'next-themes';
import { InteractionErrors } from './InteractionErrors';
import { MarkdownRenderer } from './MarkdownRenderer';
import { InteractionFooter } from './InteractionFooter';
import { aiRelayApi } from '~/api/bloefish/ai-relay';
import { FormatDuration } from '~/components/atoms/FormatDuration';
import { friendlyAiRelayOptions } from '~/utils/ai-providers';

interface InteractionContentProps {
	interaction: Interaction;
}

export const InteractionContent: React.FC<InteractionContentProps> = ({
	interaction,
}) => {
	const { data: supportedModels } = aiRelayApi.useListSupportedQuery();
	const theme = useTheme();
	const onlyErrors = Boolean(interaction.messageContent === '' && interaction.errors?.length > 0)
	const isBot = interaction.owner.type === 'bot';

	return (
		<Stack gap={2} overflow={'scroll'} maxW={'100%'}>
			{onlyErrors ? (
				<InteractionErrors errors={interaction.errors} />
			) : (
				<Card.Root
					position={'relative'}
					borderRadius={'lg'}
					blur={'10px'}
					width={'fit-content'}
					maxW={'100%'}
					minW={'200px'}
					background={theme.resolvedTheme === 'dark' ? 'rgb(17 17 17 / 40%)' : 'rgb(255 255 255 / 40%)'}
				>
					<Card.Body px={isBot ? 6 : 4} py={isBot ? 6 : 3}>
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

					{isBot && (
						<React.Fragment>
							<Badge
								top={'-1px'}
								right={'-1px'}
								borderTop={0}
								borderTopLeftRadius={0}
								borderBottomRightRadius={0}
								position={'absolute'}
								variant={'surface'}
								size={'xs'}
								colorPalette={'pink'}
								px={2}
								py={1}
							>
								{friendlyAiRelayOptions(interaction.aiRelayOptions, supportedModels?.models)}
							</Badge>
							
							<Badge
								bottom={'-1px'}
								right={'-1px'}
								borderBottom={0}
								borderBottomLeftRadius={0}
								borderTopRightRadius={0}
								position={'absolute'}
								variant={'surface'}
								size={'xs'}
								colorPalette={'gray'}
								px={2}
								py={1}
							>
								<FormatDuration start={interaction.createdAt} />
							</Badge>
						</React.Fragment>
					)}
				</Card.Root>
			)}
			
			<InteractionFooter interaction={interaction} />
		</Stack>
	)
};
