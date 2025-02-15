import { combineReducers } from '@reduxjs/toolkit';
import { userApi } from './user';

export const bloefishApiReducers = {
	[userApi.reducerPath]: userApi.reducer,
};

export const bloefishApiMiddleware = [
	userApi.middleware,
];
