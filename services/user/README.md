# svc_user

## Description

This service is a simple user service that is used to manage users.

## Base URL

`http://svc_user.bloefish.local:4001/`

## CRPC - RPC transport

### Base URL

`http://svc_user.bloefish.local:4001/rpc/<version>/<endpoint>`

### Versions

- `2025-02-12` - Initial version

### Endpoints

#### `get_user_by_id`

Gets a user by their user ID.

**Contract**

```typescript
interface Request {
	user_id: string;
}

interface Response {
	user: {
		id: string;
		default_user: boolean;
		created_at: string;
		updated_at: string | null;
		deleted_at: string | null;
	}
}
```

#### `get_or_create_default_user`

Get the default user for the platform. If the user does not exist, it will be created.

**Contract**

```typescript
type Request = null;

interface Response {
	user: {
		id: string;
		default_user: boolean;
		created_at: string;
		updated_at: string | null;
		deleted_at: string | null;
	}
}
```
