# svc_ai_relay

## Description

This service

## Base URL

`http://svc_ai_relay.bloefish.local:4002/`

### RPC transport

### Base URL

`http://svc_ai_relay.bloefish.local:4002/rpc/<version>/<endpoint>`

### Versions

- `2025-02-12` - Initial version

### Endpoints

#### `list_supported`

Lists the supported AI providers and models.

**Contract**

```typescript
type Request = null;

interface Response {
	providers: {
		id: string;
		name: string;
		models: {
			id: string;
			description: string;
		}[];
	}[];
}
```

#### `invoke_conversation_message`

This will invoke a call to an AI model, passing in a full conversation, wait for the response, and return it.

**Contract**

```typescript
interface Request {
	owner: {
		type: 'user';
		identifier: string;
	};
	messages: {
		content: string;
		owner: {
			type: 'user' | 'bot';
			identifier: string;
		};
		file_ids: string[];
	}[];
	ai_relay_options: {
		provider_id: 'open_ai';
		model_id: string;
	};
}

interface Response {
	message_response: string;
}
```

#### `invoke_streaming_conversation_message`

This will invoke a call to an AI model, passing in a full conversation, and stream the response back via the provided streaming channel id. The full response will be returned in the response of the request.

**Contract**

```typescript
interface Request {
	streaming_channel_id: string;
	owner: {
		type: 'user';
		identifier: string;
	};
	messages: {
		content: string;
		owner: {
			type: 'user' | 'bot';
			identifier: string;
		};
		file_ids: string[];
	}[];
	ai_relay_options: {
		provider_id: 'open_ai';
		model_id: string;
	};
}

interface Response {
	message_response: string;
}
```
