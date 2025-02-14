import { Avatar, Box, HStack, Separator, Stack, Text } from "@chakra-ui/react";
import { SidebarButton } from "./components/molecules/SidebarButton";
import { LuBox, LuFolderOpen, LuMessageCirclePlus, LuSettings2, LuWorkflow } from "react-icons/lu";

export const Sidebar: React.FC = () => {
	return (
		<Box
			width="250px"
			color="white"
			padding={5}
			paddingRight={2}
		>
			<Stack>
				<HStack gap={3}>
					<Avatar.Root variant={"subtle"} size={"xs"} fontSize={"md"}>
						🐡
					</Avatar.Root>
					<Text fontSize={"md"} fontWeight={"bolder"} color={'WindowText'}>
						Bloefish
					</Text>
				</HStack>
				<Separator
					marginY={2}
					marginRight={-2}
					marginLeft={-5}
				/>
				<Stack gap={2} marginRight={-2}>
					<SidebarButton content={"New conversation"} icon={<LuMessageCirclePlus />} />
					<SidebarButton content={"Conversations"} icon={<LuFolderOpen />} />
					<SidebarButton content={"Files"} icon={<LuBox />} />
					<SidebarButton content={"Workflows"} icon={<LuWorkflow />} />
					<SidebarButton content={"Preferences"} icon={<LuSettings2 />} />
				</Stack>
			</Stack>
		</Box>
	);
};
