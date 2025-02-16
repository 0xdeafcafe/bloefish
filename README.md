# bloefish

## Introduction

Bloefish is a local AI framework, aiming to provide a simple way to interact with various
AI models and providers, via a few simple to use touch points.

<p align="center">
	<img height="700" src="./assets/web.png">
</p>

## Services

| Service        | Port   | Description                                                         | Readme                            |
| -------------- | ------ | ------------------------------------------------------------------- | --------------------------------- |
| `ai_relay`     | `4003` | Handles relaying requests to various AI models and providers.       | [View](./services/airelay/README.md) |
| `conversation` | `4002` | Handles the creation and management of interactions (messages).     | [View](./services/conversation/README.md) |
| `user`         | `4001` | Handles the creation and management of users.                       | [View](./services/user/README.md) |
| `stream`       | `4004` | Handles streaming.                                                  | [View](./services/stream/README.md) |
| `file_upload`   | `4005` | Handles management, including uploading, of files.                   | [View](./services/fileupload/README.md) |

## Applications

| Application | Port   | Description                                                         | Readme                            |
| ----------- | ------ | ------------------------------------------------------------------- | --------------------------------- |
| `web`       | `5000` | A web application for interacting with Bloefish's backend.           | [View](./applications/web/README.md) |
| `cli`	      | `N/A`  | A command line interface for interacting with Bloefish's backend.    | [View](./applications/cli/README.md) |

Bloefish also exposes a full RPC API, which Bloefish itself uses, which can be used. Each
backend service exposes a full readme, and an API definition can be found [here](./beak).

## Usage

### Requirements

- Docker
- Golang
- Node.js

To run the project, fire the following commands into your terminal:

```
# This is temporary until I think of a better solution
$ echo "127.0.0.1 app.bloefish.local" | sudo tee -a /etc/hosts
$ echo "127.0.0.1 svc_ai_relay.bloefish.local" | sudo tee -a /etc/hosts
$ echo "127.0.0.1 svc_conversation.bloefish.local" | sudo tee -a /etc/hosts
$ echo "127.0.0.1 svc_file_upload.bloefish.local" | sudo tee -a /etc/hosts
$ echo "127.0.0.1 svc_stream.bloefish.local" | sudo tee -a /etc/hosts
$ echo "127.0.0.1 svc_user.bloefish.local" | sudo tee -a /etc/hosts
$ echo "127.0.0.1 storageminio" | sudo tee -a /etc/hosts

# Get the project going
$ make install
$ make build
$ make start
$ open http://app.bloefish.local:4169/testing
```

## Current annoyances

- Each service needs a hosts entry
	- `echo "127.0.0.1 svc_xxx.bloefish.local" | sudo tee -a /etc/hosts`
- The storage utility (minio), needs a special hosts entry
	- `echo "127.0.0.1 storageminio" | sudo tee -a /etc/hosts`
