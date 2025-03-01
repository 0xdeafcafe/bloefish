import { useCallback } from 'react';
import { useAppDispatch, useAppSelector } from '~/store';
import { fileUploadApi } from '~/api/bloefish/file-upload';
import {
	addFileUpload,
	updateFileUploadStatus,
	updateFileUploadProgress,
	setFileUploadFileId,
	removeFileUpload,
} from '../store';
import type { FileMetadata } from '../store/types';
import type { Actor } from '~/api/bloefish/shared.types';
import { generateRandomString } from '~/utils/random';

export const useFileUpload = (identifier: string) => {
	const dispatch = useAppDispatch();
	const chatInput = useAppSelector((state) => state.chatInput[identifier]);
	const files = chatInput?.files || {};

	const [createUploadTrigger] = fileUploadApi.useCreateUploadMutation();
	const [confirmUploadTrigger] = fileUploadApi.useConfirmUploadMutation();

	const uploadFile = useCallback(async (file: File, owner: Actor) => {
		const fileMetadata: FileMetadata = {
			name: file.name,
			type: file.type,
			size: file.size,
			lastModified: file.lastModified,
		};

		const fileUploadId = generateRandomString(10);

		dispatch(addFileUpload({
			identifier,
			fileUploadId,
			fileMetadata,
		}));

		try {
			dispatch(updateFileUploadStatus({ identifier, fileUploadId, status: 'uploading' }));

			const createUploadResponse = await createUploadTrigger({
				name: file.name,
				size: file.size,
				mimeType: file.type,
				owner,
			}).unwrap();

			dispatch(setFileUploadFileId({
				identifier,
				fileUploadId,
				fileId: createUploadResponse.id,
				uploadUrl: createUploadResponse.uploadUrl
			}));

			const uploadResponse = await uploadToPresignedUrl(
				createUploadResponse.uploadUrl,
				file,
				(progress) => {
					dispatch(updateFileUploadProgress({ identifier, fileUploadId, progress }));
				}
			);

			if (!uploadResponse.ok) {
				throw new Error(`Failed to upload file: ${uploadResponse.statusText}`);
			}

			dispatch(updateFileUploadStatus({ identifier, fileUploadId, status: 'confirming' }));

			await confirmUploadTrigger({
				fileId: createUploadResponse.id,
			}).unwrap();

			dispatch(updateFileUploadStatus({ identifier, fileUploadId, status: 'ready' }));
		} catch (error) {
			console.error('Upload failed:', error);
			dispatch(updateFileUploadStatus({
				identifier,
				fileUploadId,
				status: 'error',
				error: error instanceof Error ? error.message : 'Unknown error',
			}));
		}
	}, [dispatch, identifier, files, createUploadTrigger, confirmUploadTrigger]);

	const remove = useCallback((fileUploadId: string) => {
		dispatch(removeFileUpload({ identifier, fileUploadId }));
	}, [dispatch, identifier]);

	return {
		files,
		uploadFile,
		removeFileUpload: remove,
	};
};

async function uploadToPresignedUrl(
	presignedUrl: string,
	file: File,
	onProgress?: (progress: number) => void
): Promise<Response> {
	return new Promise((resolve, reject) => {
		const xhr = new XMLHttpRequest();

		xhr.upload.addEventListener('progress', (event) => {
			if (event.lengthComputable && onProgress) {
				const progress = Math.round((event.loaded / event.total) * 100);
				onProgress(progress);
			}
		});

		xhr.addEventListener('load', () => {
			if (xhr.status >= 200 && xhr.status < 300) {
				resolve(new Response(null, {
					status: xhr.status,
					statusText: xhr.statusText,
				}));
			} else {
				reject(new Error(`HTTP Error: ${xhr.status}`));
			}
		});

		xhr.addEventListener('error', () => {
			reject(new Error('Network error occurred'));
		});

		xhr.addEventListener('abort', () => {
			reject(new Error('Upload aborted'));
		});

		xhr.open('PUT', presignedUrl);
		xhr.setRequestHeader('Content-Type', file.type);
		xhr.send(file);
	});
}
