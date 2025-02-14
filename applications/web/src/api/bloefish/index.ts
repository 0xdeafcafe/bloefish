import { combineReducers } from '@reduxjs/toolkit';
import { user } from './user';

export const bloefishApiReducer = combineReducers({
	[user.reducerPath]: user.reducer,
});

export const bloefishApiMiddleware = [
	user.middleware,
];
