package controllers

import (
	"crm-server/models"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type RoleController struct {
	beego.Controller
}

// 角色列表
func (c *RoleController) RoleList() {
	var roleList []*models.Role
	err := db.Cols("id", "name", "desc").Find(&roleList)
	if err != nil {
		msg := err.Error()
		log.Error(msg)
		c.Redirect("/error?msg="+msg, 302)
		return
	}
	c.Data["roleList"] = roleList
	c.TplName = "role/role_list.html"
}

// 添加角色页面
func (c *RoleController) AddPage() {
	c.TplName = "role/add_role.html"
}

// 添加角色提交
func (c *RoleController) Add() {
	moduleId := c.GetString("module_id", "")
	powerList := c.GetString("power_list", "")
	name := c.GetString("name")
	desc := c.GetString("desc", "")

	// 开启事务
	dbSession := db.NewSession()
	defer dbSession.Close()
	dbSession.Begin()

	// 添加角色
	var role = models.Role{
		Name:    name,
		Desc:    desc,
		Modules: moduleId,
	}
	_, err := dbSession.Insert(&role)
	if err != nil {
		dbSession.Rollback();
		log.Error(err.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err.Error(),
		}
		c.ServeJSON()
		return
	}
	// 给该角色添加权限
	powerSlice := strings.Split(powerList, ",")
	for _, v := range powerSlice {
		enforcer.AddPolicy(strconv.Itoa(role.Id), v, "*")
	}
	// 提交事务
	dbSession.Commit()

	c.Data["json"] = map[string]interface{}{
		"code": 200,
		"msg":  "添加成功",
	}
	c.ServeJSON()
}

// 角色详情
func (c *RoleController) Detail() {
	id, _ := c.GetInt("id")
	var role models.Role
	_, err := db.Id(id).Get(&role)
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	// 查询权限列表
	var powerMap []map[string]string
	err2 := db.Table("casbin_rule").Cols("v1").Where("v0=?", id).Find(&powerMap)
	if err2 != nil {
		log.Error(err2.Error())
		c.Redirect("/error?msg="+err2.Error(), 302)
		return
	}
	// 转为切片
	var powerList []string
	for _, v := range powerMap {
		powerList = append(powerList, v["v1"])
	}

	c.Data["powerList"] = strings.Join(powerList, ",")
	c.Data["role"] = role
	c.TplName = "role/role_detail.html"
}

// 角色提交修改
func (c *RoleController) Mod() {
	id, _ := c.GetInt("id")
	moduleId := c.GetString("module_id", "")
	powerList := c.GetString("power_list", "")
	name := c.GetString("name")
	desc := c.GetString("desc", "")

	// 开启事务
	dbSession := db.NewSession()
	defer dbSession.Close()

	dbSession.Begin()
	// 角色信息
	var role = models.Role{
		Name:    name,
		Desc:    desc,
		Modules: moduleId,
	}
	_, err := dbSession.ID(id).Update(&role)
	if err != nil {
		dbSession.Rollback()
		log.Error(err.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err.Error(),
		}
		c.ServeJSON()
		return
	}
	// 查询该角色之前的权限并删除
	var ruleList []map[string]string
	err2 := dbSession.Table("casbin_rule").Cols("v0", "v1", "v2").Where("v0=?", id).Find(&ruleList);
	if err2 != nil {
		dbSession.Rollback()
		log.Error(err2.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err2.Error(),
		}
		c.ServeJSON()
		return
	}

	for _, v := range ruleList {
		enforcer.RemovePolicy(v["v0"], v["v1"], v["v2"])
	}

	// 给该角色添加权限
	powerSlice := strings.Split(powerList, ",")
	for _, v := range powerSlice {
		enforcer.AddPolicy(strconv.Itoa(id), v, "*")
	}
	// 提交事务
	dbSession.Commit()

	c.Data["json"] = map[string]interface{}{
		"code": 200,
		"msg":  "保存成功",
	}
	c.ServeJSON()
}

// 删除角色
func (c *RoleController) Del() {
	id,_ := c.GetInt("id")
	dbSession := db.NewSession()
	defer dbSession.Close()
	dbSession.Begin()

	_,err := dbSession.ID(id).Delete(new(models.Role))
	if err != nil {
		dbSession.Rollback()
		log.Error(err.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 将该角色的用户的角色字段改为0
	var user = models.Users {
		RoleId: 0,
	}

	_,err2 := dbSession.Where("role_id=?",id).Cols("role_id").Update(&user)
	if err2 != nil {
		dbSession.Rollback()
		log.Error(err2.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg": err2.Error(),
		}
		c.ServeJSON()
		return
	}

	// 删除对应的权限
	var ruleList []map[string]string
	err3 := dbSession.Table("casbin_rule").Cols("v0", "v1", "v2").Where("v0=?", id).Find(&ruleList);
	if err3 != nil {
		dbSession.Rollback()
		log.Error(err3.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err3.Error(),
		}
		c.ServeJSON()
		return
	}

	for _, v := range ruleList {
		enforcer.RemovePolicy(v["v0"], v["v1"], v["v2"])
	}

	dbSession.Commit()

	c.Data["json"] = map[string]interface{}{
		"code": 200,
		"msg": "删除成功",
	}
	c.ServeJSON()
}
