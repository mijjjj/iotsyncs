package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"hotgo/internal/library/hggen/internal/service"
)

var (
	Install = cInstall{}
)

type cInstall struct {
	g.Meta `name:"install" brief:"install gf binary to system (might need root/admin permission)"`
}

type cInstallInput struct {
	g.Meta `name:"install"`
}

type cInstallOutput struct{}

func (c cInstall) Index(ctx context.Context, in cInstallInput) (out *cInstallOutput, err error) {
	err = service.Install.Run(ctx)
	return
}
