import { Center, EmptyState, Icon, VStack } from "@chakra-ui/react";
import { Helmet } from "react-helmet-async";
import { LuFishOff } from "react-icons/lu";

interface NotFoundProps {
	Icon?: React.FC;
	title?: string;
	pageTitle?: string;
	description?: string;
}

export const NotFound: React.FC<NotFoundProps> = ({
	Icon,
	title,
	pageTitle,
	description,
}) => (
	<Center height={'full'}>
		<Helmet>
			<title>{`${pageTitle ?? 'Not Found'} | Bloefish`}</title>
		</Helmet>

		<EmptyState.Root size={'lg'}>
			<EmptyState.Content>
				<EmptyState.Indicator>
					{Icon ? (
						<Icon />
					) : (
						<LuFishOff style={{
							'transform': 'scaleY(-1);',
						}} />
					)}
				</EmptyState.Indicator>
				<VStack textAlign="center">
					<EmptyState.Title>
						{title ?? 'That page doesn\'t exist...'}
					</EmptyState.Title>
					<EmptyState.Description>
						{description ?? 'Have you tried looking harder?'}
					</EmptyState.Description>
				</VStack>
			</EmptyState.Content>
		</EmptyState.Root>
	</Center>
);
