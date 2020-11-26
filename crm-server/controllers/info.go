package controllers

import "github.com/astaxie/beego"

type InfoController struct {
	beego.Controller
}

func (c *InfoController) ErrorPage(){
	msg := c.GetString("msg", "")

	c.Data["msg"] = msg
	c.TplName = "info/error.html"
}