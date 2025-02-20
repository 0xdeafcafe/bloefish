import { useEffect } from 'react';
import { useAppDispatch } from '~/store';
import type { StreamMessage } from './stream.types';
import { addInteractionError, addInteractionFragment, updateConversationTitle, updateInteractionMessageContent } from '~/features/conversations/store';
import camelcaseKeys from 'camelcase-keys';

export function useStreamListener() {
	const dispatch = useAppDispatch();

	useEffect(() => {
		const url = 'ws://svc_stream.bloefish.local:4004/ws';
		let ws = new WebSocket(url);

		ws.onmessage = (event) => {
			const message = camelcaseKeys(JSON.parse(event.data)) as StreamMessage;

			switch (message.type) {
				case 'message_fragment': {
					const [conversationId, interactionId] = message.channelId.split('/');
					if (!conversationId || !interactionId) return;

					if (interactionId === 'title') {
						dispatch(updateConversationTitle({
							conversationId,
							title: message.messageFragment,
							treatAsFragment: true,
						}));
						break;
					}

					dispatch(addInteractionFragment({
						conversationId,
						interactionId,
						fragment: message.messageFragment,
					}));
					break;
				}

				case 'message_full': {
					const [conversationId, interactionId] = message.channelId.split('/');
					if (!conversationId || !interactionId) return;

					if (interactionId === 'title') {
						dispatch(updateConversationTitle({
							conversationId,
							title: message.messageFull,
							treatAsFragment: false,
						}));
						break;
					}

					dispatch(updateInteractionMessageContent({
						conversationId,
						interactionId,
						content: message.messageFull,
					}));
					break;
				}

				case 'error': {
					const [conversationId, interactionId] = message.channelId.split('/');
					if (!conversationId || !interactionId) return;

					dispatch(addInteractionError({
						conversationId,
						interactionId,
						error: message.error,
					}));
					break;
				}

				default:
					break;
			}
		};

		ws.onopen = () => {
			console.log('Connected to stream');
		};

		ws.onclose = () => {
			console.log('Disconnected from stream');

			setTimeout(() => {
				ws = new WebSocket(url);
			}, 1000);
		};

		ws.onerror = (error) => {
			console.error('WebSocket error:', error);
		};

		return () => {
			ws.close();
		};
	}, []);
}
