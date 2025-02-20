import { Card, Flex } from "@chakra-ui/react";
import { motion } from "motion/react";
import type React from "react";

const Root: React.FC<React.PropsWithChildren> = ({ children }) => (
	<MotionCardRoot
		flex={'1'}
		variant={'outline'}
		boxShadow={'md'}
		p={0}
		overflow={'hidden'}
		display={'grid'}
		gridTemplateRows={'auto fit-content(1fr)'}

		initial={{ scale: 0.95, opacity: 0, y: 10 }}
		animate={{ scale: 1, opacity: 1, y: 0 }}
		exit={{ scale: 0.95, opacity: 0, y: 10 }}
		transition={{ duration: 0.3 }}
	>
		{children}
	</MotionCardRoot>
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

const Body: React.FC<React.PropsWithChildren> = ({ children }) => (
	<Card.Body
		p={0}
		height={'100%'}
		overflowY={'scroll'}
	>
		{children}
	</Card.Body>
);

export const Panel = {
	Root,
	Header,
	Body,
};

const MotionCardRoot = motion.create(Card.Root);
