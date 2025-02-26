import { createApi } from '@reduxjs/toolkit/query/react';
import { createBaseQueryWithSnake } from './base';
import type { CreateSkillSetRequest, ListSkillSetsByOwnerRequest, ListSkillSetsByOwnerResponse } from './skill-set.types';

export const skillSetApi = createApi({
	reducerPath: 'api.bloefish.skill_set',
	baseQuery: createBaseQueryWithSnake('http://svc_skill_set.bloefish.local:4006/rpc/'),

	endpoints: (builder) => ({
		createSkillSet: builder.mutation<void, CreateSkillSetRequest>({
			query: (body) => ({
				url: '2025-02-12/create_skill_set',
				body,
			}),
			// async onQueryStarted(req, { dispatch, queryFulfilled }) {
			// 	await queryFulfilled;

			// 	dispatch(skillSetApi.endpoints.createSkillSet.initiate())
			// },
		}),

		listSkillSetsByOwner: builder.query<ListSkillSetsByOwnerResponse, ListSkillSetsByOwnerRequest>({
			query: (body) => ({
				url: '2025-02-12/list_skill_sets_by_owner',
				body,
			}),
		}),
	}),
});
