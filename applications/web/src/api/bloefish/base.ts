import { fetchBaseQuery, type BaseQueryFn, type FetchArgs } from '@reduxjs/toolkit/query/react';
import snakecaseKeys from 'snakecase-keys';
import camelcaseKeys from 'camelcase-keys';

export const createBaseQueryWithSnake = (baseUrl: string): BaseQueryFn<string | FetchArgs, unknown, unknown> => {
	const baseQuery = fetchBaseQuery({
		baseUrl,
		method: 'POST',
	});

	return async (args, api, extraOptions) => {
		// If there's a body, convert its keys to `snake_case`
		if (args && typeof args === 'object' && 'body' in args && args.body) {
			(args as FetchArgs).body = snakecaseKeys(args.body, { deep: true });
		}

		const result = await baseQuery(args, api, extraOptions);

		// If the response has data, convert its keys to `camelCase`
		if (result.data) {
			result.data = camelcaseKeys(result.data as Record<string, unknown>, { deep: true });
		}

		return result;
	};
};
