import { createApi } from '@reduxjs/toolkit/query/react';
import { createBaseQueryWithSnake } from './base';
import type { ConfirmUploadRequest, CreateUploadRequest, CreateUploadResponse } from './file-upload.types';

export const fileUploadApi = createApi({
	reducerPath: 'api.bloefish.file-upload',
	baseQuery: createBaseQueryWithSnake('http://svc_file_upload.bloefish.local:4005/rpc/'),

	endpoints: (builder) => ({
		createUpload: builder.mutation<CreateUploadResponse, CreateUploadRequest>({
			query: (body) => ({
				url: '2025-02-12/create_upload',
				body,
			}),
		}),

		confirmUpload: builder.mutation<void, ConfirmUploadRequest>({
			query: (body) => ({
				url: '2025-02-12/confirm_upload',
				body,
			}),
		}),
	}),
});
