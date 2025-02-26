import { Text, type TextProps } from '@chakra-ui/react';
import {
	format,
	formatDistance,
	addMinutes,
	addHours,
	addDays,
	differenceInMilliseconds,
} from 'date-fns';
import { Tooltip } from '../ui/tooltip';
import { useState, useEffect } from 'react';

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
	const [now, setNow] = useState(new Date());
	const startDateTime = new Date(start);
	const endDateTime = end ? new Date(end) : now;

	useEffect(() => {
		const getNextUpdateTime = (from: Date, to: Date) => {
			const seconds = Math.floor(differenceInMilliseconds(from, to) / 1000);

			if (seconds < 3600) {
				return addMinutes(from, 1); // Next minute mark
			}
			if (seconds < 86400) {
				return addHours(from, 1); // Next hour mark
			}

			return addDays(from, 1); // Next day mark
		};

		const nextUpdate = getNextUpdateTime(startDateTime, now);
		const timeUntilUpdate = differenceInMilliseconds(nextUpdate, now);
		
		if (timeUntilUpdate <= 0) return;

		const timer = setTimeout(() => setNow(new Date()), timeUntilUpdate);
		return () => clearTimeout(timer);
	}, [now, startDateTime]);

	const distance = formatDistance(startDateTime, endDateTime, { addSuffix: !hideSuffix });

	return (
		<Tooltip content={format(endDateTime, 'PPpp')}>
			<Text textDecoration={'dotted'} textDecorationLine={'underline'} cursor={'help'} {...rest}>
				{distance}
			</Text>
		</Tooltip>
	);
};
