import { Blockquote, Card, Link, List, Table, Text } from '@chakra-ui/react';
import { styled } from 'styled-components';
import { LuExternalLink } from 'react-icons/lu'
import remarkGfm from 'remark-gfm';
import Markdown from 'react-markdown';
import { useTheme } from 'next-themes';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { twilight, coy } from 'react-syntax-highlighter/dist/esm/styles/prism';

interface MarkdownRendererProps {
	markdown: string;
}

export const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({
	markdown,
}) => {
	const theme = useTheme();

	return (
		<Text fontSize={'sm'} as={'div'}>
			<MarkdownWrapper>
				<Markdown
					disallowedElements={[]}
					remarkPlugins={[remarkGfm]}
					components={{
						code(props) {
							const { children, className, node, ...rest } = props
							const match = /language-(\w+)/.exec(className || '')
							return match ? (
								// @ts-expect-error too lazy to fix this
								<SyntaxHighlighter
									{...rest}
									PreTag={'div'}
									children={String(children).replace(/\n$/, '')}
									language={match[1]}
									style={theme.resolvedTheme === 'dark' ? twilight : coy}
								/>
							) : (
								<code {...rest} className={className}>
									{children}
								</code>
							)
						},

						h1: (props) => <Text as="h1" fontSize="2xl" fontWeight="bold" {...props} />,
						h2: (props) => <Text as="h2" fontSize="xl" fontWeight="bold" {...props} />,
						h3: (props) => <Text as="h3" fontSize="lg" fontWeight="bold" {...props} />,
						h4: (props) => <Text as="h4" fontSize="md" fontWeight="bold" {...props} />,
						h5: (props) => <Text as="h5" fontSize="sm" fontWeight="bold" {...props} />,
						h6: (props) => <Text as="h6" fontSize="xs" fontWeight="bold" {...props} />,

						a: (props) => {
							return (
								<Link colorPalette={'purple'} whiteSpace={'wrap'} variant={'underline'} {...props}>
									{props.children} <LuExternalLink />
								</Link>
							)
						},

						blockquote: (props) => (
							<Blockquote.Root>
								<Blockquote.Content>
									{props.children}
								</Blockquote.Content>
							</Blockquote.Root>
						),

						ul: (props) => (
							<List.Root as={'ul'} ml={4} mt={-3}>
								{props.children}
							</List.Root>
						),
						ol: (props) => (
							<List.Root as={'ol'} ml={4} mt={-3}>
								{props.children}
							</List.Root>
						),
						li: (props) => (<List.Item>{props.children}</List.Item>),

						table: (props) => (
							<Card.Root overflow={'hidden'}>
								<Table.ScrollArea>
									<Table.Root size={'sm'} variant={'outline'} striped stickyHeader interactive>
										{props.children}
									</Table.Root>
								</Table.ScrollArea>
							</Card.Root>
						),
						tr: (props) => <Table.Row>{props.children}</Table.Row>,
						th: (props) => <Table.ColumnHeader>{props.children}</Table.ColumnHeader>,
						td: (props) => <Table.Cell>{props.children}</Table.Cell>,
						tbody: (props) => <Table.Body>{props.children}</Table.Body>,
						thead: (props) => <Table.Header>{props.children}</Table.Header>,
						tfoot: (props) => <Table.Footer>{props.children}</Table.Footer>,
					}}
				>
					{markdown}
				</Markdown>
			</MarkdownWrapper>
		</Text>
	);
};

const MarkdownWrapper = styled.div`
	& > * {
		margin-bottom: 16px;
	}

	& > *:last-child {
		margin-bottom: 0;
	}
`;
