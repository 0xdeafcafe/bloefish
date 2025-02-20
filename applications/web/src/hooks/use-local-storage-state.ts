import { useState } from 'react';

type LocalStorageKey = 'chat_input.selected_model';

export function useLocalStorageState<T>(
	key: LocalStorageKey,
	initialValue?: T
): [T | undefined, (value: T | undefined) => void] {
	const [storedValue, setStoredValue] = useState<T>(() => {
		try {
			const item = window.localStorage.getItem(key);

			return item ? JSON.parse(item) : initialValue;
		} catch (error) {
			console.error(error);
			return initialValue;
		}
	});

	const setValue = (value: T) => {
		try {
			setStoredValue(value);
			window.localStorage.setItem(key, JSON.stringify(value));
		} catch (error) {
			console.error(error);
		}
	};

	return [storedValue, setValue];
}
