import { Flex, Spinner } from '@chakra-ui/react';
import { userApi } from '~/api/bloefish/user';

export const EnsureReadiness: React.FC<React.PropsWithChildren> = ({ children }) => {
	const {
		isError,
		isFetching,
		isLoading,
		isUninitialized,
	} = userApi.useGetOrCreateDefaultUserQuery();

	if (isFetching || isLoading || isUninitialized) {
		// render loading
		return (
			<Flex w={'full'} h={'full'} justify={'center'} align={'center'}>
				<Spinner />
			</Flex>
		);
	}
	if (isError) {
		// render try again
		return 'error';
	}

	return children;
};
