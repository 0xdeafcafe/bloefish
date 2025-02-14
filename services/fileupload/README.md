# svc_file_upload

## Description

This service handles the management of files, including uploading and downloading them.

## Base URL

`http://svc_file_upload.bloefish.local:4005/`

## RPC transport

### Base URL

`http://svc_file_upload.bloefish.local:4005/rpc/<version>/<endpoint>`

### Versions

- `2025-02-12` - Initial version

### Endpoints

#### `create_upload`

Creates a new file upload, which will create a file object and return a URL to upload the file to.

**Contract**

```typescript
interface Request {
	name: string;
	size: number;
	mime_type: string;
	owner: {
		type: 'user';
		identifier: string;
	};
}

interface Response {
	id: string;
	upload_url: string;
}
```

#### `confirm_upload`

Confirms that the file has been uploaded, until this is called an uploaded file can not be used. This will also validate that the uploaded file matches the data that was provided in the `create_upload` endpoint.

**Contract**

```typescript
interface Request {
	file_id: string;
}

type Response = null;
```

#### `get_file`

Gets a file object by its ID. If `include_access_url` is set to `true`, then the `presigned_access_url` will be included in the response.

If `access_url_expiry_seconds` is set to a number, then the `presigned_access_url` will expire after that many seconds. If it is set to `null`, then the URL will have a default expiry time of 15 minutes.

**Contract**

```typescript
interface Request {
	file_id: string;
	include_access_url: boolean;
	access_url_expiry_seconds: number | null;
}

interface Response {
	id: string;
	name: string;
	size: number;
	mime_type: string;
	owner: {
		type: 'user';
		identifier: string;
	};
	presigned_access_url: string | null;
}
```


## Uploading to an upload url

```bash
```bash
curl -X PUT \
	-T "file-text.txt" \
	"http://storageminio.bloefish.local:9000/bloefish-svc-files/file_000000D1ldWfjVInywfVDtwm4bBqb?X-Amz-Algorithm=AWS4-HMAC-SHA256..."
```
