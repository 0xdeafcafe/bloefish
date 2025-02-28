import type { Actor } from './shared.types';

export interface CreateUploadRequest {
	name: string;
	size: number;
	mime_type: string;
	owner: Actor;
}

export interface CreateUploadResponse {
	id: string;
	upload_url: string;
}

export interface ConfirmUploadRequest {
	file_id: string;
}
