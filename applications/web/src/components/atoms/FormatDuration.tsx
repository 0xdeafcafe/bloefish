import { Text, type TextProps } from '@chakra-ui/react';
import { format, formatDistance }from'date-fns';
import { Tooltip } from '../ui/tooltip';

interface FormatDurationProps extends TextProps {
	start: string;
	end?: string;
	hideSuffix?: boolean;
}

export const FormatDuration: React.FC<FormatDurationProps> = ({
	start,
	end,
	hideSuffix,
	...rest
}) => {
	const startDateTime = new Date(start);
	const endDateTime = end ? new Date(end) : new Date();
	const distance = formatDistance(startDateTime, endDateTime, { addSuffix: !hideSuffix });

	return (
		<Tooltip content={format(endDateTime, 'PPpp')}>
			<Text textDecoration={'dotted'} textDecorationLine={'underline'} cursor={'help'} {...rest}>
				{distance}
			</Text>
		</Tooltip>
	);
};
