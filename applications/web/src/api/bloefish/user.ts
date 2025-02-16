import { createApi } from '@reduxjs/toolkit/query/react';
import type { GetOrCreateDefaultUserResponse } from './user.types';
import { createBaseQueryWithSnake } from './base';

export const userApi = createApi({
	reducerPath: 'api.bloefish.user',
	baseQuery: createBaseQueryWithSnake('http://svc_user.bloefish.local:4001/rpc/'),

	endpoints: (builder) => ({
		getOrCreateDefaultUser: builder.query<GetOrCreateDefaultUserResponse, void>({
			query: () => `2025-02-12/get_or_create_default_user`,
		}),
	}),
})
