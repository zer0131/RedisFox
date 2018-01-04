package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (mc *MainController) Get() {
	mc.Data["Website"] = "beego.me"
	mc.Data["Email"] = "astaxie@gmail.com"
	mc.TplName = "index.tpl"
}
