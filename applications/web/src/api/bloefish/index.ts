import { aiRelayApi } from './ai-relay';
import { conversationApi } from './conversation';
import { userApi } from './user';

export const bloefishApiReducers = {
	[aiRelayApi.reducerPath]: aiRelayApi.reducer,
	[conversationApi.reducerPath]: conversationApi.reducer,
	[userApi.reducerPath]: userApi.reducer,
};

export const bloefishApiMiddleware = [
	aiRelayApi.middleware,
	conversationApi.middleware,
	userApi.middleware,
];
