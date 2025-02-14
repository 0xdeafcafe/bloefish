import { combineReducers, configureStore } from '@reduxjs/toolkit';
import { apiMiddleware, apiReducer } from './api';

export type RootState = ReturnType<typeof store.getState>;

export const store = configureStore({
	reducer: combineReducers({
		api: apiReducer,
	}),
	// Adding the api middleware enables caching, invalidation, polling and other features of RTK Query
	middleware: (getDefaultMiddleware) => 
		getDefaultMiddleware().concat(
			apiMiddleware,
		),
});
