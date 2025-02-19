import { Box, HStack, Icon, Text } from "@chakra-ui/react";
import { motion } from "motion/react";
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
		<MotionLink 
			to={path}
			initial={{ scale: 1 }}
			whileTap={{ scale: 0.95 }}
		>
			<Box
				background={active ? 'bg.emphasized' : 'transparent'}
				paddingX={3}
				paddingY={'5px'}
				borderWidth={'1px'}
				borderColor={active ? 'border.emphasized' : 'transparent'}
				borderRadius={'md'}
				onClick={onClick}
				cursor={'pointer'}
				zIndex={10000}
			>
				<HStack>
					<Icon color={active ? 'MenuText' : 'GrayText'}>
						{icon}
					</Icon>
					<Text fontSize={"sm"} color={'MenuText'} fontWeight={'semibold'}>{content}</Text>
				</HStack>
			</Box>
		</MotionLink>
	);
};

const MotionLink = motion(Link);
