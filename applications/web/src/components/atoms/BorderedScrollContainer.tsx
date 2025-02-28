import { Box, type BoxProps } from '@chakra-ui/react';
import React, { useRef, useState } from 'react';

interface BorderedScrollContainerProps extends BoxProps {
	triggerOffset: number;
}

export const BorderedScrollContainer: React.FC<React.PropsWithChildren<BorderedScrollContainerProps>> = ({
	children,
	triggerOffset,
	...rest
}) => {

	const scrollContainerRef = useRef<HTMLDivElement>(null);
	const [showBorder, setShowBorder] = useState(false);

	return (
		<Box
			{...rest}
			overflow={'auto'}
			ref={scrollContainerRef}
			onScroll={() => {
				const scrollTop = scrollContainerRef.current?.scrollTop ?? 0;
				setShowBorder(scrollTop > triggerOffset);
			}}
			borderTopWidth={showBorder ? '1px' : '0'}
			borderTopColor={'border.emphasized'}
			boxShadow={showBorder ? 'inner' : 'none'}
		>
			{children}
		</Box>
	);
};
