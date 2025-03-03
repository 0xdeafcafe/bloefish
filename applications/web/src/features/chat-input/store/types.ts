import type { EnrichedAiModel } from '~/api/bloefish/ai-relay.types';

export type UploadStatus = 'pending' | 'uploading' | 'confirming' | 'ready' | 'error';

export interface FileMetadata {
	name: string;
	type: string;
	size: number;
	lastModified?: number;
}

export interface ChatInputPlugin {
	identifier: string;
}

export interface ChatInputState {
	prompt: string;
	skillSetIds: string[];
	destinationModel?: EnrichedAiModel;
	files: Record<string, FileUploadState>;
}

export interface FileUploadState {
	fileMetadata: FileMetadata;
	status: UploadStatus;
	progress?: number;
	error?: string;
	fileId?: string;
	uploadUrl?: string;
}

export interface UpdatePromptPayload extends ChatInputPlugin {
	prompt: string;
}

export interface UpdateAiRelayOptionsPayload extends ChatInputPlugin {
	destinationModel: EnrichedAiModel;
}

export interface UpdateSkillSetIdsPayload extends ChatInputPlugin {
	skillSetIds: string[];
}

export interface AddFileUploadPayload extends ChatInputPlugin {
	fileUploadId: string;
	fileMetadata: FileMetadata;
}

export interface UpdateFileUploadStatusPayload extends ChatInputPlugin {
	fileUploadId: string;
	status: UploadStatus;
	error?: string;
}

export interface UpdateFileUploadProgressPayload extends ChatInputPlugin {
	fileUploadId: string;
	progress: number;
}

export interface SetFileUploadFileIdPayload extends ChatInputPlugin {
	fileUploadId: string;
	fileId: string;
	uploadUrl: string;
}

export interface RemoveFileUploadPayload extends ChatInputPlugin {
	fileUploadId: string;
}
