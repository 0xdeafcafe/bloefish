import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

export interface OmniBarState {
	open: boolean;
	query: string;
}

export const omniBarSlice = createSlice({
	name: 'conversations',
	initialState: {
		open: false,
		query: '',
	} as OmniBarState,
	reducers: {
		openOmni: (state, { payload }: PayloadAction<{ keepQuery?: boolean } | void>) => {
			state.open = true;
			if (!payload?.keepQuery) state.query = '';
		},
		closeOmni: (state, { payload }: PayloadAction<{ keepQuery?: boolean } | void>) => {
			state.open = false;
			if (!payload?.keepQuery) state.query = '';
		},
		toggleOmni: (state, { payload }: PayloadAction<{ keepQuery?: boolean } | void>) => {
			state.open = !state.open;
			if (!payload?.keepQuery) state.query = '';
		},
		updateOmniQuery: (state, { payload }: PayloadAction<string>) => {
			state.query = payload;
		},
	},
});

export const {
	closeOmni,
	openOmni,
	toggleOmni,
	updateOmniQuery,
} = omniBarSlice.actions;
export const omniBarReducer = omniBarSlice.reducer;
