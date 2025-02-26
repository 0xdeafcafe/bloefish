# svc_conversation

## Description

This service

## Base URL

`http://svc_conversation.bloefish.local:4002/`

## RPC transport

### Base URL

`http://svc_conversation.bloefish.local:4002/rpc/<version>/<endpoint>`

### Versions

- `2025-02-12` - Initial version

### Endpoints

#### `create_conversation`

Creates a new conversation.

**Contract**

```typescript
interface Request {
	idempotency_key: string;
	owner: {
		type: 'user';
		identifier: string;
	};
	ai_relay_options: {
		provider_id: 'open_ai';
		model_id: string;
	};
}

interface Response {
	conversation_id: string;
	owner: {
		type: 'user';
		identifier: string;
	};
	ai_relay_options: {
		provider_id: 'open_ai';
		model_id: string;
	};
	
	title: string | null;
	stream_channel_id: string;

	created_at: string; // ISO 8601
	updated_at: string; // ISO 8601
	deleted_at: string | null; // ISO 8601
}
```

#### `create_conversation_message`

Creates a new message in a conversation. This message will be appended to the conversation, and the entire conversation chain will be sent to the AI relay.

If `ai_relay_options` is set to `null`, then the `ai_relay_options` set on the conversation will be used.

**Contract**

```typescript
interface Request {
	conversation_id: string;
	idempotency_key: string;
	message_content: string;
	file_ids: string[];
	skill_set_ids: string[];

	owner: {
		type: 'user';
		identifier: string;
	};
	ai_relay_options: {
		provider_id: 'open_ai';
		model_id: string;
	} | null;

	options: {
		use_streaming: boolean;
	};
}

interface Response {
	conversation_id: string;

	input_interaction: {
		id: string;
		file_ids: string[];
		skill_set_ids: string[];
	
		marked_as_excluded_at: string | null; // ISO 8601

		message_content: string;
		errors: {
			code: string;
			message: string;
			reasons: {
				code: string;
				message: string;
			}[];
		}[];

		owner: {
			type: 'user';
			identifier: string;
		};
		ai_relay_options: {
			provider_id: 'open_ai';
			model_id: string;
		};

		created_at: string; // ISO 8601
		updated_at: string; // ISO 8601
		deleted_at: string | null; // ISO 8601
		completed_at: string | null; // ISO 8601
	};
	response_interaction: {
		id: string;
		file_ids: string[];
		skill_set_ids: string[];
	
		marked_as_excluded_at: string | null; // ISO 8601

		message_content: string;
		errors: {
			code: string;
			message: string;
			reasons: {
				code: string;
				message: string;
			}[];
		}[];

		owner: {
			type: 'user';
			identifier: string;
		};
		ai_relay_options: {
			provider_id: 'open_ai';
			model_id: string;
		};

		created_at: string; // ISO 8601
		updated_at: string; // ISO 8601
		deleted_at: string | null; // ISO 8601
		completed_at: string | null; // ISO 8601
	};

	stream_channel_id: string;
}
```

#### `get_conversation_with_interactions`

Gets a conversation with all of its interactions.

**Contract**

```typescript
interface Request {
	conversation_id: string;
}

interface Response {
	conversation_id: string;
	owner: {
		type: 'user';
		identifier: string;
	};
	ai_relay_options: {
		provider_id: 'open_ai';
		model_id: string;
	};

	title: string | null;
	stream_channel_id: string;

	interactions: {
		id: string;
		file_ids: string[];
		skill_set_ids: string[];

		marked_as_excluded_at: string | null; // ISO 8601

		message_content: string;
		errors: {
			code: string;
			message: string;
			reasons: {
				code: string;
				message: string;
			}[];
		}[];

		owner: {
			type: 'user' | 'bot';
			identifier: string;
		};
		ai_relay_options: {
			provider_id: 'open_ai';
			model_id: string;
		};

		created_at: string; // ISO 8601
		updated_at: string; // ISO 8601
		deleted_at: string | null; // ISO 8601
		completed_at: string | null; // ISO 8601
	}[];

	created_at: string; // ISO 8601
	updated_at: string; // ISO 8601
	deleted_at: string | null; // ISO 8601
}
```

#### `list_conversations_with_interactions`

Lists all conversations with all of their interactions.

**Contract**

```typescript
interface Request {
	owner: {
		type: 'user';
		identifier: string;
	};
}

interface Response {
	conversations: {
		id: string;
		owner: {
			type: 'user';
			identifier: string;
		};
		ai_relay_options: {
			provider_id: 'open_ai';
			model_id: string;
		};

		title: string | null;
		stream_channel_id: string;

		interactions: {
			id: string;
			file_ids: string[];
			skill_set_ids: string[];

			marked_as_excluded_at: string | null; // ISO 8601

			message_content: string;
			errors: {
				code: string;
				message: string;
				reasons: {
					code: string;
					message: string;
				}[];
			}[];

			owner: {
				type: 'user' | 'bot';
				identifier: string;
			};
			ai_relay_options: {
				provider_id: 'open_ai';
				model_id: string;
			};

			created_at: string; // ISO 8601
			updated_at: string; // ISO 8601
			deleted_at: string | null; // ISO 8601
			completed_at: string | null; // ISO 8601
		}[];

		created_at: string; // ISO 8601
		updated_at: string; // ISO 8601
		deleted_at: string | null; // ISO 8601
	}[];
}
```

#### `delete_conversations`

Deletes conversations and their interactions by the conversation ID.

**Contract**

```typescript
interface Request {
	conversation_ids: string[];
}

type Response = null;
```

#### `delete_interactions`

Deletes interactions by the their ID.

**Contract**

```typescript
interface Request {
	interaction_ids: string;
}

type Response = null;
```

#### `update_interaction_excluded_state`

Updates the excluded state of an interaction.

**Contract**

```typescript
interface Request {
	interaction_id: string;
	excluded: boolean;
}

type Response = null;
```
