export type LocalStorageKey =
	'chat_input.selected_model' |
	`chat_input.selected_model.${string}`;

export function getValue<T>(key: LocalStorageKey, initialValue?: T): T | undefined {
	try {
		const item = window.localStorage.getItem(key);

		return item ? JSON.parse(item) : initialValue;
	} catch (error) {
		console.error(error);

		return initialValue;
	}
}

export function setValue<T>(key: LocalStorageKey, value: T | undefined): void {
	try {
		window.localStorage.setItem(key, JSON.stringify(value));
	} catch (error) {
		console.error(error);
	}
}
