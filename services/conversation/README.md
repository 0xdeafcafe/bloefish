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

	// If streaming is used in the conversation at all, this is the channel the messages will be attached to.
	stream_channel_id: string;
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
	interaction_id: string;

	response_interaction_id: string;

	// If streaming is used in the conversation at all, this is the channel the messages will be attached to.
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
	id: string;
	owner: {
		type: 'user';
		identifier: string;
	};
	ai_relay_options: {
		provider_id: 'open_ai';
		model_id: string;
	};
	interactions: {
		id: string;
		owner: {
			type: 'user' | 'bot';
			identifier: string;
		};
		message_content: string;
		file_ids: string[];
		ai_relay_options: {
			provider_id: 'open_ai';
			model_id: string;
		} | null;
		created_at: string; // ISO 8601
	}[];
	created_at: string; // ISO 8601
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
		interactions: {
			id: string;
			owner: {
				type: 'user' | 'bot';
				identifier: string;
			};
			message_content: string;
			file_ids: string[];
			ai_relay_options: {
				provider_id: 'open_ai';
				model_id: string;
			} | null;
			created_at: string; // ISO 8601
		}[];
		created_at: string; // ISO 8601
	}[];
}
```
