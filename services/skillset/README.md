# svc_skill_set

## Description

This service handles the management of files, including uploading and downloading them.

## Base URL

`http://svc_skill_set.bloefish.local:4006/`

## RPC transport

### Base URL

`http://svc_skill_set.bloefish.local:4006/rpc/<version>/<endpoint>`

### Versions

- `2025-02-12` - Initial version

### Endpoints

#### `create_skill_set`

Creates a new skill set.

**Contract**

```typescript
interface Request {
	name: string;
	icon: string;
	description: string;
	prompt: string;

	owner: {
		type: 'user';
		identifier: string;
	};
}

type Response = null;
```

#### `update_skill_set`

Creates a new skill set.

**Contract**

```typescript
interface Request {
	skill_set_id: string;
	name: string;
	icon: string;
	description: string;
	prompt: string;
}

type Response = null;
```

#### `delete_skill_set`

Creates a new skill set.

**Contract**

```typescript
interface Request {
	skill_set_id: string;
}

type Response = null;
```

#### `get_skill_set`

Gets a skill set by ID.

**Contract**

```typescript
interface Request {
	id: string;
	name: string;
	icon: string;
	description: string;
	prompt: string;

	owner: {
		type: 'user';
		identifier: string;
	};

	created_at: string;
	updated_at: string;
	deleted_at: string | null;
}

type Response = null;
```

#### `list_skill_sets_by_owner`

Lists all skill sets for an owner.

**Contract**

```typescript
interface Request {
	owner: {
		type: 'user';
		identifier: string;
	};
}

interface Response {
	skill_sets: {
		id: string;
		name: string;
		icon: string;
		description: string;
		prompt: string;

		owner: {
			type: 'user';
			identifier: string;
		};

		created_at: string;
		updated_at: string;
		deleted_at: string | null;
	}[];
}
```
