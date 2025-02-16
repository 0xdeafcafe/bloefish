export interface StreamMessageFull {
	channelId: string;
	type: 'message_full';
	messageFull: string;
	messageFragment: null;
}

export interface StreamMessageFragment {
	channelId: string;
	type: 'message_fragment';
	messageFull: null;
	messageFragment: string;
}

export type StreamMessage = StreamMessageFull | StreamMessageFragment;
