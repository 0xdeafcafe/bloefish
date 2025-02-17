# svc_stream

## Description

This service is responsible for handling streaming.

## Base URL

`http://svc_stream.bloefish.local:4004/`

## RPC transport

### Base URL

`http://svc_stream.bloefish.local:4004/rpc/<version>/<endpoint>`

### Versions

- `2025-02-12` - Initial version

### Endpoints

#### `send_message_full`

Sends a full message.

**Contract**

```typescript
interface Request {
	channel_id: string;
	message_content: string;
}

type Response = null;
```

#### `send_message_fragment`

Sends a message fragment.

**Contract**

```typescript
interface Request {
	channel_id: string;
	message_content: string;
}

type Response = null;
```

#### `send_error`

Sends an error message.

**Contract**

```typescript
interface Request {
	channel_id: string;
	error: {
		code: number;
		meta: Record<string, unknown>;
		reasons: { code: string, meta: Record<string, unknown>, reasons: unknown[] }[];
	};
}

type Response = null;
```

## WebSocket transport

### Base URL

`http://svc_stream.bloefish.local:4004/ws`

### Message Structure

```typescript
interface Message {
	channel_id: string;
	message_id: string;
	type: 'message_full' | 'message_fragment' | 'error_message';
	message_full: string | null; // Only set if type is 'message_full'
	message_fragment: string | null; // Only set if type is 'message_fragment'
	error: {
		code: number;
		meta: Record<string, unknown>;
		reasons: { code: string, meta: Record<string, unknown>, reasons: unknown[] }[];
	} | null; // Only set if type is 'error'
}
```
