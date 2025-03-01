import type { Actor } from './shared.types';

export interface CreateUploadRequest {
	name: string;
	size: number;
	mimeType: string;
	owner: Actor;
}

export interface CreateUploadResponse {
	id: string;
	uploadUrl: string;
}

export interface ConfirmUploadRequest {
	fileId: string;
}
