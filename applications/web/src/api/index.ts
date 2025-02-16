import { bloefishApiMiddleware, bloefishApiReducers } from './bloefish';

export const apiReducers = {
	...bloefishApiReducers,
};

export const apiMiddleware = [
	...bloefishApiMiddleware,
];
