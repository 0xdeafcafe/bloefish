package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/libraries/contexts"

	serviceAIRelay "github.com/0xdeafcafe/bloefish/services/airelay/cmd"
	serviceConversation "github.com/0xdeafcafe/bloefish/services/conversation/cmd"
	serviceFileUpload "github.com/0xdeafcafe/bloefish/services/fileupload/cmd"
	serviceStream "github.com/0xdeafcafe/bloefish/services/stream/cmd"
	serviceUser "github.com/0xdeafcafe/bloefish/services/user/cmd"
)

type ServiceBoot func(ctx context.Context, args []string) error

var (
	serviceDefinitions = map[string]ServiceBoot{
		"ai_relay":     serviceAIRelay.Root,
		"conversation": serviceConversation.Root,
		"file_upload":  serviceFileUpload.Root,
		"stream":       serviceStream.Root,
		"user":         serviceUser.Root,
	}
)

func main() {
	os.Exit(root(os.Args))
}

func root(args []string) int {
	ctx := context.Background()

	args = args[1:]

	if err := handle(ctx, args); err != nil {
		clog.Get(ctx).WithError(err).Error("service exited unexpectedly")
		return 1
	}

	return 0
}

func handle(ctx context.Context, args []string) error {
	defer clog.HandlePanic(context.Background(), false)

	if len(args) == 0 {
		return errors.New("no command specified")
	}

	cmd := args[0]
	args = args[1:]

	ctx = contexts.SetServiceInfo(ctx, contexts.ServiceInfo{
		Environment:     config.GetEnvironmentName(),
		Service:         cmd,
		ServiceHTTPName: "svc_" + cmd,
		GitRepository:   "https://github.com/0xdeafcafe/bloefish",
	})

	fn, ok := serviceDefinitions[cmd]
	if !ok {
		return fmt.Errorf("unrecognized command: %s", cmd)
	}

	return fn(ctx, args)
}
