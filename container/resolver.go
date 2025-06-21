package container

import (
	"context"
	"github.com/web-page-analysis/bootstrap"
)

type Container struct {
	OBAdapter OutBoundConnection
}

func Resolver(ctx context.Context,
	conf bootstrap.Config) *Container {
	//outbound connection resolver
	outBoundConnectionAdapter := InitOutBoundConnection(conf)

	return &Container{
		OBAdapter: outBoundConnectionAdapter,
	}
}
