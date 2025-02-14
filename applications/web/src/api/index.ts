import { combineReducers } from '@reduxjs/toolkit';
import { bloefishApiMiddleware, bloefishApiReducer } from './bloefish';

export const apiReducer = combineReducers({
	bloefish: bloefishApiReducer,
});

export const apiMiddleware = [
	...bloefishApiMiddleware,
];
