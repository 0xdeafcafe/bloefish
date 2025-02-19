import { useEffect } from 'react';
import { useAppDispatch } from '~/store';
import type { StreamMessage } from './stream.types';
import { addInteractionFragment, updateInteractionMessageContent } from '~/features/conversations/store';
import camelcaseKeys from 'camelcase-keys';

export function useStreamListener() {
	const dispatch = useAppDispatch();

	useEffect(() => {
		const url = 'ws://svc_stream.bloefish.local:4004/ws';
		let ws = new WebSocket(url);

		ws.onmessage = (event) => {
			const data = camelcaseKeys(JSON.parse(event.data)) as StreamMessage;

			switch (data.type) {
				case 'message_fragment': {
					const [conversationId, interactionId] = data.channelId.split('/');
					if (!conversationId || !interactionId) return;

					dispatch(addInteractionFragment({
						conversationId,
						interactionId,
						fragment: data.messageFragment,
					}));
					break;
				}

				case 'message_full':
					const [conversationId, interactionId] = data.channelId.split('/');
					if (!conversationId || !interactionId) return;

					dispatch(updateInteractionMessageContent({
						conversationId,
						interactionId,
						content: data.messageFull,
					}));
					break;

				default:
					break;
			}

			if (data.type !== 'message_fragment') return;

			
		};

		ws.onopen = () => {
			console.log('Connected to stream');
		};

		ws.onclose = () => {
			console.log('Disconnected from stream');

			setTimeout(() => {
				ws = new WebSocket(url);
			}, 500);
		};

		ws.onerror = (error) => {
			console.error('WebSocket error:', error);
		};

		return () => {
			ws.close();
		};
	}, []);
}
