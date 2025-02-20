import { Alert, Box, Stack, Text } from '@chakra-ui/react';
import { motion } from 'motion/react';
import { useState } from 'react';
import { LuFileQuestion } from 'react-icons/lu';
import type { Cher } from '~/api/bloefish/shared.types';
import { Button } from '~/components/ui/button';
import { generateRandomString } from '~/utils/random';

interface InteractionErrorsProps {
	errors: Cher[] | null;
}

export const InteractionErrors: React.FC<InteractionErrorsProps> = ({
	errors
}) => {
	if (!errors?.length) return null;

	return (
		<Stack>
			{errors.map((error) => (
				<InteractionError error={error} key={generateRandomString(5)} />
			))}
		</Stack>
	);
};

const InteractionError: React.FC<{ error: Cher }> = ({ error }) => {
	const [expanded, setExpanded] = useState(false);

	switch (error.code) {
		case 'ai_model_not_found':
			return (
				<Alert.Root status={'error'} variant={'subtle'}>
					<Alert.Indicator>
						<LuFileQuestion />
					</Alert.Indicator>
					<Alert.Content overflowX={'scroll'}>
						<Alert.Title>
							Model not supported
						</Alert.Title>
						<Alert.Description>
							Please select another model from the model selector under the chat bar.
						</Alert.Description>
					</Alert.Content>
				</Alert.Root>
			);

		default:
			return (
				<Alert.Root status={'error'} variant={'subtle'}>
					<Alert.Indicator />
					<Alert.Content overflowX={'scroll'}>
						<Alert.Title>
							Stacked it
						</Alert.Title>
						<Alert.Description>
							There was an unknown system issue. Show the error 
							<MotionBox
								mt={2}
								maxW={'full'}
								initial={false}
								animate={{
									height: expanded ? 'auto' : 0,
									opacity: expanded ? 1 : 0
								}}
								transition={{ duration: 0.2 }}
								// overflow="hidden"
							>
								<Text textStyle={'xs'}>
									<pre>
										{JSON.stringify(error, null, 2)}
									</pre>
								</Text>
							</MotionBox>
						</Alert.Description>
					</Alert.Content>
					<Button
						alignSelf={'flex-start'}
						variant={'outline'}
						size={'sm'}
						fontWeight={'medium'}
						onClick={() => setExpanded(prev => !prev)}
					>
						{expanded ? 'Collapse error' : 'Expand error'}
					</Button>
				</Alert.Root>
			);
	}
}

const MotionBox = motion.create(Box);
