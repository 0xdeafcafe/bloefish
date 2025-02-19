import { Button, Flex, Icon, Link as ChakraLink, Text } from '@chakra-ui/react';
import type React from 'react';
import { Link } from 'react-router';

interface OmniButtonProps extends React.PropsWithChildren {
	iconElement: React.ReactElement;
	onClick: () => void;
}

export const OmniButton: React.FC<OmniButtonProps> = ({
	iconElement,
	children,
	onClick,
}) => (
	<Button
		colorPalette={'purple'}
		variant={'ghost'}
		size={'xs'}
		justifyContent={'flex-start'}
		onClick={onClick}
	>
		<Flex gap={2} overflow={'hidden'}>
			<Icon>
				{iconElement}
			</Icon>
			<Text textOverflow={'ellipsis'} overflow={'hidden'} whiteSpace={'nowrap'}>
				{children}
			</Text>
		</Flex>
	</Button>
);

interface OmniLinkProps extends React.PropsWithChildren {
	iconElement: React.ReactElement;
	href: string;
}

// sinful code below
export const OmniLink: React.FC<OmniLinkProps> = ({
	iconElement,
	children,
	href,
}) => (
	<ChakraLink asChild overflow={'hidden'} colorPalette={'purple'}>
		<Link
			to={href}
			style={{ textDecoration: 'none'}}
		>
			<Button
				variant={'ghost'}
				justifyContent={'flex-start'}
				size={'xs'}
				width={'full'}
				tabIndex={-1}
			>
				<Flex gap={2} overflow={'hidden'}>
					<Icon>
						{iconElement}
					</Icon>
					<Text textOverflow={'ellipsis'} overflow={'hidden'} whiteSpace={'nowrap'}>
						{children}
					</Text>
				</Flex>
			</Button>
		</Link>
	</ChakraLink>
);
