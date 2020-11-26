package controllers

import (
	"crm-server/models"
	"crm-server/utils"
	"github.com/astaxie/beego"
	"strings"
)

type IndexController struct {
	beego.Controller
}

// 初始化日志引擎
var log = utils.Logger
// 初始化db引擎
var db = models.DB
// 初始化casbin Enforcer
var enforcer = models.Enforcer

// 首页
func (c *IndexController) Index() {
	c.TplName = "index/index.html"
}

// 头部模块
func (c *IndexController) Header() {
	// 获取用户名
	userId := c.GetSession("user_id")

	var user models.Users
	_, err := db.Cols("name").ID(userId).Get(&user)
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	c.Data["user"] = user
	c.TplName = "index/header.html"
}

// 头部模块
func (c *IndexController) MainArea() {
	c.TplName = "index/main.html"
}

// 菜单模块
func (c *IndexController) Menu() {
	// 获取用户id
	userId := c.GetSession("user_id")
	var user models.Users
	_, err := db.Id(userId).Get(&user)
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	// 查询角色id
	var role models.Role
	_,err2 := db.Cols("id", "modules").ID(user.RoleId).Get(&role)
	if err2 != nil {
		log.Error(err2.Error())
		c.Redirect("/error?msg="+err2.Error(), 302)
		return
	}
	// 查询权限列表
	var ruleList []map[string]string
	err3 := db.Table("casbin_rule").Cols("v1").ID(role.Id).Find(&ruleList)
	if err3 != nil {
		log.Error(err3.Error())
		c.Redirect("/error?msg="+err3.Error(), 302)
		return
	}
	var powerList []string
	for _,v := range ruleList {
		powerList = append(powerList, v["v1"])
	}
	c.Data["powerList"] = powerList
	c.Data["moduleList"] = strings.Split(role.Modules,",")
	c.Data["userId"] = userId
	c.TplName = "index/menu.html"
}

// 欢迎页
func (c *IndexController) Welcome() {
	c.TplName = "index/welcome.html"
}
