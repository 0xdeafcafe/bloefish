import { Card, Flex } from "@chakra-ui/react";
import type React from "react";

const Body: React.FC<React.PropsWithChildren> = ({ children }) => (
	<Card.Body
		p={0}
		maxH={'calc(100vh - 3.5rem)'}
		minH={'700px'}
		overflow={'hidden'}
	>
		{children}
	</Card.Body>
);

const Header: React.FC<React.PropsWithChildren> = ({ children }) => (
	<Card.Header
		py={0}
		px={4}
		height={'60px'}
		borderBottomWidth={'1px'}
		borderBottomStyle={'solid'}
		borderBottomColor={'border.emphasized'}
	>
		<Flex
			justify={'space-between'}
			align={'center'}
			height={'full'}
		>
			{children}
		</Flex>
	</Card.Header>
);

export const Panel = {
	Body,
	Header,	
};
