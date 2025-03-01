import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import { getValue, setValue } from '~/utils/localstorage';
import type {
	ChatInputPlugin,
	ChatInputState,
	UpdateAiRelayOptionsPayload,
	UpdatePromptPayload,
	UpdateSkillSetIdsPayload,
	AddFileUploadPayload,
	UpdateFileUploadStatusPayload,
	UpdateFileUploadProgressPayload,
	RemoveFileUploadPayload,
	SetFileUploadFileIdPayload,
} from './types';

const initialState: Record<string, ChatInputState | undefined> = {};
const initialChatInputState = (identifier: string): ChatInputState => ({
	skillSetIds: [],
	prompt: '',
	destinationModel: getValue(`chat_input.selected_model.${identifier}`) ?? getValue('chat_input.selected_model'),
	files: {},
});

export const chatInputSlice = createSlice({
	name: 'chatInput',
	initialState: initialState,
	reducers: {
		initializeChatInput: (state, { payload }: PayloadAction<ChatInputPlugin>) => {
			state[payload.identifier] = initialChatInputState(payload.identifier);
		},
		updatePrompt: (state, { payload }: PayloadAction<UpdatePromptPayload>) => {
			let chatInput = state[payload.identifier];
			if (!chatInput) chatInput = state[payload.identifier] = initialChatInputState(payload.identifier);

			chatInput.prompt = payload.prompt;
		},
		updateDestinationModel: (state, { payload }: PayloadAction<UpdateAiRelayOptionsPayload>) => {
			let chatInput = state[payload.identifier];
			if (!chatInput) chatInput = state[payload.identifier] = initialChatInputState(payload.identifier);

			chatInput.destinationModel = payload.destinationModel;
			setValue(`chat_input.selected_model.${payload.identifier}`, payload.destinationModel);
		},
		updateSkillSetIds: (state, { payload }: PayloadAction<UpdateSkillSetIdsPayload>) => {
			let chatInput = state[payload.identifier];
			if (!chatInput) chatInput = state[payload.identifier] = initialChatInputState(payload.identifier);

			chatInput.skillSetIds = payload.skillSetIds;
		},

		addFileUpload: (state, { payload }: PayloadAction<AddFileUploadPayload>) => {
			let chatInput = state[payload.identifier];
			if (!chatInput) chatInput = state[payload.identifier] = initialChatInputState(payload.identifier);

			if (!chatInput.files) chatInput.files = {};

			chatInput.files[payload.fileUploadId] = {
				fileMetadata: payload.fileMetadata,
				status: 'pending',
			};

			return state;
		},

		updateFileUploadStatus: (state, { payload }: PayloadAction<UpdateFileUploadStatusPayload>) => {
			const chatInput = state[payload.identifier];
			if (!chatInput?.files) return;

			const upload = chatInput.files[payload.fileUploadId];
			if (upload) {
				upload.status = payload.status;
				if (payload.error !== void 0) {
					upload.error = payload.error;
				}
			}
		},

		updateFileUploadProgress: (state, { payload }: PayloadAction<UpdateFileUploadProgressPayload>) => {
			const chatInput = state[payload.identifier];
			if (!chatInput?.files) return;

			const upload = chatInput.files[payload.fileUploadId];
			if (upload) {
				upload.progress = payload.progress;
			}
		},

		setFileUploadFileId: (state, { payload }: PayloadAction<SetFileUploadFileIdPayload>) => {
			const chatInput = state[payload.identifier];
			if (!chatInput?.files) return;

			const upload = chatInput.files[payload.fileUploadId];
			if (upload) {
				upload.fileId = payload.fileId;
				upload.uploadUrl = payload.uploadUrl;
			}
		},

		removeFileUpload: (state, { payload }: PayloadAction<RemoveFileUploadPayload>) => {
			const chatInput = state[payload.identifier];
			if (!chatInput?.files) return;

			delete chatInput.files[payload.fileUploadId];
		},
	},
});

export const {
	initializeChatInput,
	updatePrompt,
	updateDestinationModel,
	updateSkillSetIds,
	addFileUpload,
	updateFileUploadStatus,
	updateFileUploadProgress,
	setFileUploadFileId,
	removeFileUpload,
} = chatInputSlice.actions;
export const chatInputReducer = chatInputSlice.reducer;
