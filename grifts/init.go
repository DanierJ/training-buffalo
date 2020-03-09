package grifts

import (
	"github.com/DanierJ/div_manager/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
