import { Card, Text } from "@chakra-ui/react";

interface HeaderCardProps {
	title: string;
	description: string;
}

export const HeaderCard: React.FC<HeaderCardProps> = ({
	title,
	description,
}) => (
	<Card.Root
		variant={'elevated'}
		overflow={'hidden'}
		background={'bg.muted'}
		backgroundSize={'100vw 500px'}
	>
		<Card.Body p={6}>
			<Text textStyle={'4xl'} fontWeight={'bolder'} textShadow={'2xl'}>{title}</Text>
			<Text textStyle={'sm'} textShadow={'md'}>{description}</Text>
		</Card.Body>
	</Card.Root>
);
