import { Box, HStack, Icon, Text } from "@chakra-ui/react";
import { Link } from "react-router";

interface SidebarButtonProps {
	active: boolean;
	icon: React.ReactElement;
	path: string;
	content: string;
	onClick?: () => void;
}

export const SidebarLink: React.FC<SidebarButtonProps> = ({
	active,
	icon,
	path,
	content,
	onClick,
}) => {
	return (
		<Link to={path}>
			<Box
				background={active ? 'bg.emphasized' : 'transparent'}
				paddingX={3}
				paddingY={'5px'}
				borderRadius={5}
				borderRightRadius={0}
				onClick={onClick}
				cursor={'pointer'}
			>
				<HStack>
					<Icon color={active ? 'MenuText' : 'GrayText'}>
						{icon}
					</Icon>
					<Text fontSize={"sm"} color={'MenuText'} fontWeight={'bold'}>{content}</Text>
				</HStack>
			</Box>
		</Link>
	);
};
