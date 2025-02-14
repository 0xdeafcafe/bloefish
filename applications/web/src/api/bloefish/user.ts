import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type { GetOrCreateDefaultUserResponse } from './user.types';

export const user = createApi({
	reducerPath: 'user',
	baseQuery: fetchBaseQuery({
		baseUrl: 'http://svc_user.bloefish.local:4001/rpc/',
		method: 'POST',
	}),

	endpoints: (builder) => ({
	  getOrCreateDefaultUser: builder.query<GetOrCreateDefaultUserResponse, void>({
		query: () => `2025-02-12/get_or_create_default_user`,
	  }),
	}),
  })
  