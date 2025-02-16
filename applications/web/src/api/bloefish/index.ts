import { conversationApi } from './conversation';
import { userApi } from './user';

export const bloefishApiReducers = {
	[conversationApi.reducerPath]: conversationApi.reducer,
	[userApi.reducerPath]: userApi.reducer,
};

export const bloefishApiMiddleware = [
	conversationApi.middleware,
	userApi.middleware,
];
