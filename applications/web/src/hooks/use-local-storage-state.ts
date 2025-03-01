import { useState } from 'react';
import { getValue, setValue as setValueUtil, type LocalStorageKey } from '~/utils/localstorage';

export function useLocalStorageState<T>(
	key: LocalStorageKey,
	initialValue?: T
): [T | undefined, (value: T | undefined) => void] {
	const [storedValue, setStoredValue] = useState<T | undefined>(() => {
		return getValue(key, initialValue);
	});

	const setValue = (value: T | undefined) => {
		setValueUtil(key, value);
		setStoredValue(value);
	};

	return [storedValue, setValue];
}
