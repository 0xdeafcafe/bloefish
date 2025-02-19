import { combineReducers, configureStore } from '@reduxjs/toolkit';
import { apiMiddleware, apiReducers } from './api';
import { useDispatch, useSelector } from 'react-redux';
import { conversationsReducer } from './features/conversations/store';
import { omniBarReducer } from './features/omnibar/store';

export type RootState = ReturnType<typeof store.getState>;

export const store = configureStore({
	reducer: combineReducers({
		...apiReducers,
		conversations: conversationsReducer,
		omniBar: omniBarReducer,
	}),
	devTools: true,
	middleware: (getDefaultMiddleware) => 
		getDefaultMiddleware().concat(
			apiMiddleware,
		),
});

export type AppDispatch = typeof store.dispatch;
export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector = useSelector.withTypes<RootState>();
