import { createApi } from '@reduxjs/toolkit/query/react';
import { createBaseQueryWithSnake } from './base';
import {  } from '~/features/conversations/store';
import type { ListSupportedResponse } from './ai-relay.types';

export const aiRelayApi = createApi({
	reducerPath: 'api.bloefish.aiRelay',
	baseQuery: createBaseQueryWithSnake('http://svc_ai_relay.bloefish.local:4003/rpc/'),

	endpoints: (builder) => ({
		listSupported: builder.query<ListSupportedResponse, void>({
			query: (body) => ({
				url: '2025-02-12/list_supported',
				body,
			}),
		}),
	}),
});
