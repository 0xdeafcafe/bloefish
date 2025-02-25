import { Box, Card, ClipboardRoot, Code } from '@chakra-ui/react';
import { useTheme } from 'next-themes';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { twilight, coy } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { ClipboardIconButton } from '~/components/ui/clipboard';

interface CodeViewerProps {
	className: string | undefined;
	children: React.ReactNode | undefined;
}

export const CodeViewer: React.FC<CodeViewerProps> = ({
	children,
	className,
}) => {
	const theme = useTheme();
	const match = /language-(\w+)/.exec(className || '');

	if (match) {
		return (
			<Card.Root borderRadius={'lg'}>
				<Card.Body p={0} position={'relative'}>
					<SyntaxHighlighter
						PreTag={'div'}
						// eslint-disable-next-line react/no-children-prop
						children={String(children).replace(/\n$/, '')}
						language={match[1]}
						style={theme.resolvedTheme === 'dark' ? twilight : coy}
					/>
					<Box position={'absolute'} top={4} right={4}>
						<ClipboardRoot value={String(children)}>
							<ClipboardIconButton />
						</ClipboardRoot>
					</Box>
				</Card.Body>
			</Card.Root>
		)
	}

	return (
		<Code variant={'surface'}>
			{children}
		</Code>
	);
};
