import { Stack, Text } from '@chakra-ui/react';
import { format, formatDistanceToNow, formatISO }from'date-fns';
import { Tooltip } from '../ui/tooltip';

interface FormatDurationProps {
	pointInTime: string;
}

export const FormatDuration: React.FC<FormatDurationProps> = ({ pointInTime }) => {
	const pointInTimeDate = new Date(pointInTime);
	const distance = formatDistanceToNow(pointInTimeDate);

	return (
		<Tooltip content={format(pointInTimeDate, 'PPpp')}>
			<Text textDecoration={'dotted'} textDecorationLine={'underline'} cursor={'help'}>
				{distance}
			</Text>
		</Tooltip>
	);
};
