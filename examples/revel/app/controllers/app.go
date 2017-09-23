package controllers

import (
	"github.com/revel/revel"
	"fmt"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Hello(name string) revel.Result {
	return c.RenderText(fmt.Sprintf("Hello %s good to see you", name))
}
