import { Stack, Text } from "@chakra-ui/react";

interface OmniGroupProps extends React.PropsWithChildren {
	title: string;
}

export const OmniGroup: React.FC<OmniGroupProps> = ({
	title,
	children,
}) => (
	<Stack>
		<Text
			color={'MenuText'}
			fontWeight={'semibold'}
			fontSize={'xs'}
		>
			{title}
		</Text>
		{children}
	</Stack>
);
