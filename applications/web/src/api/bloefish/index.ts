import { aiRelayApi } from './ai-relay';
import { conversationApi } from './conversation';
import { skillSetApi } from './skill-set';
import { userApi } from './user';

export const bloefishApiReducers = {
	[aiRelayApi.reducerPath]: aiRelayApi.reducer,
	[conversationApi.reducerPath]: conversationApi.reducer,
	[skillSetApi.reducerPath]: skillSetApi.reducer,
	[userApi.reducerPath]: userApi.reducer,
};

export const bloefishApiMiddleware = [
	aiRelayApi.middleware,
	conversationApi.middleware,
	skillSetApi.middleware,
	userApi.middleware,
];
