import { Box, Card, Em, Flex, HStack, Icon, Link, Status, Text } from "@chakra-ui/react";
import { Sidebar } from "./features/sidebar/Sidebar";
import { LuExternalLink } from "react-icons/lu";

export const App: React.FC<React.PropsWithChildren> = ({ children }) => (
	<Flex minHeight={'100vh'} background={'Background'}>
		<Sidebar />

		<Flex flex="1" marginTop={3} marginRight={3} marginBottom={3} direction={'column'}>
			<Card.Root
				flex="1"
				variant={'outline'}
				boxShadow={'md'}
				marginBottom={3}
			>
				<Card.Body gap="2">
					{children}
				</Card.Body>
			</Card.Root>
			<Box marginRight={1}>
				<HStack justify={"end"} gap={2}>
					<Status.Root colorPalette={'green'} size={'sm'}>
						<Status.Indicator />
						Online
					</Status.Root>
					<Text fontSize={'xs'} color={'GrayText'}>|</Text>
					<Link href={'https://github.com/0xdeafcafe/bloefish'} target={'_blank'} rel={'noopener'} fontSize={'xs'}>
						GitHub
						<LuExternalLink />
					</Link>
					<Text fontSize={'xs'} color={'GrayText'}>|</Text>
					<Link href={'https://github.com/0xdeafcafe/bloefish/issues'} target={'_blank'} rel={'no'} fontSize={'xs'}>
						GitHub Issues
						<LuExternalLink />
					</Link>
					<Text fontSize={'xs'} color={'GrayText'}>|</Text>
					<Text fontSize={'xs'} color={'GrayText'}>Version: <b>dev</b></Text>
					<Text fontSize={'xs'} color={'GrayText'}>|</Text>
					<Text fontSize={'xs'} color={'GrayText'}>
						Made with ❤️ in Amsterdam
					</Text>
				</HStack>
			</Box>
		</Flex>
	</Flex>
);
