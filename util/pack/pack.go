package pack

import (
	"os"

	"github.com/skotchpine/xvm/util/keyval"
)

type Ctx struct {
	Path, Version string
	Config        map[string]string
}

func Context() (ctx *Ctx, err error) {
	ctx = new(Ctx)

	ctx.Path = os.Getenv("XVM_PULL_PATH")
	ctx.Version = os.Getenv("XVM_PULL_VERSION")
	ctx.Config, err = keyval.ParseString(os.Getenv("XVM_PULL_CONFIG"))

	return
}
