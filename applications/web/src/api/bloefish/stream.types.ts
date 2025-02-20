import type { Cher } from "./shared.types";

export interface StreamMessageFull {
	channelId: string;
	type: 'message_full';
	messageFull: string;
	messageFragment: null;
	error: null;
}

export interface StreamMessageFragment {
	channelId: string;
	type: 'message_fragment';
	messageFull: null;
	messageFragment: string;
	error: null;
}

export interface StreamErrorMessage {
	channelId: string;
	type: 'error';
	messageFull: null;
	messageFragment: null;
	error: Cher;
}

export type StreamMessage = StreamMessageFull | StreamMessageFragment | StreamErrorMessage;
