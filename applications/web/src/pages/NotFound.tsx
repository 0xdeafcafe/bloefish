import { Center, EmptyState, Icon, VStack } from "@chakra-ui/react";
import { LuFishOff } from "react-icons/lu";

export const NotFound: React.FC = () => (
	<Center height={'full'}>
		<EmptyState.Root size={'lg'}>
			<EmptyState.Content>
				<EmptyState.Indicator>
					<LuFishOff style={{
						'transform': 'scaleY(-1);',
					}} />
				</EmptyState.Indicator>
				<VStack textAlign="center">
					<EmptyState.Title>That page doesn't exist...</EmptyState.Title>
					<EmptyState.Description>
						Have you tried looking harder?
					</EmptyState.Description>
				</VStack>
			</EmptyState.Content>
		</EmptyState.Root>
	</Center>
);
