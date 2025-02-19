import { useState } from "react";
import { generateRandomString } from "~/utils/random";

export function useIdempotencyKey(): [string, () => void] {
	const [idempotencyKey, setIdempotencyKey] = useState(() => generateRandomString(20));

	function generateNewKey() {
		setIdempotencyKey(generateRandomString(20));
	}

	return [idempotencyKey, generateNewKey];
}
