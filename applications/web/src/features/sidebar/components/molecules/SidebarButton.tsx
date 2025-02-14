import { Box, HStack, Icon, Text } from "@chakra-ui/react";

interface SidebarButtonProps {
	icon: React.ReactElement;
	content: string;
	onClick?: () => void;
}

export const SidebarButton: React.FC<SidebarButtonProps> = ({
	icon,
	content,
	onClick,
}) => {
	return (
		<Box paddingX={2} paddingY={1} borderRadius={5} onClick={onClick}>
			<HStack>
				<Icon>
					{icon}
				</Icon>
				<Text fontSize={"sm"} color={'MenuText'}>{content}</Text>
			</HStack>
		</Box>
	);
};
