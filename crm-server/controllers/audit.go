package controllers

import (
	"crm-server/models"
	"github.com/astaxie/beego"
	"strconv"
	"time"
)

type AuditInfo struct {
	models.Audit `xorm:"extends"`
	UserName     string `xorm:""`
}

type AuditController struct {
	beego.Controller
}

// 发布审批页面
func (c *AuditController) AddPage() {
	c.TplName = "audit/add_audit.html"
}

// 提交审批
func (c *AuditController) Add() {
	title := c.GetString("title", "")
	desc := c.GetString("title", "")
	money, err := c.GetFloat("money", 0)

	if err != nil {
		log.Error(err.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err.Error(),
		}
		c.ServeJSON()
		return
	}
	// 获取session中保存的user_id
	userId := c.GetSession("user_id")

	var audit = models.Audit{
		Title:      title,
		Money:      money,
		Desc:       desc,
		CreateTime: strconv.FormatInt(time.Now().Unix(), 10),
		UserId:     userId.(int),
	}

	_, err2 := db.Insert(audit)
	if err2 != nil {
		log.Error(err2.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err2.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 200,
		"msg":  "添加成功",
	}
	c.ServeJSON()
}

// 一级确认列表
func (c *AuditController) FirstConfirmList() {
	var auditList []*AuditInfo

	err := db.Table("audit").Alias("a").Select("a.id, a.create_time, a.money, a.title, u.name as user_name").Join("inner", "users as u", "a.user_id=u.id").Where("a.first_confirm=0 and a.second_confirm=0").OrderBy("id desc").Find(&auditList)
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	c.Data["auditList"] = auditList

	c.TplName = "audit/confirm_list.html"
}

// 二级确认列表
func (c *AuditController) SecondConfirmList() {
	var auditList []*AuditInfo

	err := db.Table("audit").Alias("a").Select("a.id, a.create_time, a.money, a.title, u.name as user_name").Join("inner", "users as u", "a.user_id=u.id").Where("a.first_confirm=1 and a.second_confirm=0").OrderBy("id desc").Find(&auditList)
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	c.Data["auditList"] = auditList

	c.TplName = "audit/confirm_list.html"
}

// 一级确认
func (c *AuditController) FirstConfirm() {
	id, _ := c.GetInt("id")
	status, _ := c.GetInt("status", 1)

	var audit = models.Audit{
		FirstConfirm:     status,
		FirstConfirmTime: strconv.FormatInt(time.Now().Unix(), 10),
	}

	row, err := db.Id(id).Where("first_confirm=0").Update(&audit)
	if err != nil {
		log.Error(err.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err.Error(),
		}
		c.ServeJSON()
		return
	}
	if row < 1 {
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "操作失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 200,
		"msg":  "操作成功",
	}
	c.ServeJSON()
}

// 二级确认
func (c *AuditController) SecondConfirm() {
	id, _ := c.GetInt("id")
	status, _ := c.GetInt("status", 1)

	var audit = models.Audit{
		SecondConfirm:     status,
		SecondConfirmTime: strconv.FormatInt(time.Now().Unix(), 10),
	}

	row, err := db.Id(id).Where("second_confirm=0").Update(&audit)
	if err != nil {
		log.Error(err.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err.Error(),
		}
		c.ServeJSON()
		return
	}
	if row < 1 {
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "操作失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 200,
		"msg":  "操作成功",
	}
	c.ServeJSON()
}

// 审批详情
func (c *AuditController) Detail() {
	id, _ := c.GetInt("id")
	// 通过id查询结款审批详情
	var auditInfo AuditInfo
	_, err := db.Table("audit").Alias("a").Select("a.id, a.create_time, a.money,a.first_confirm,a.second_confirm, a.first_confirm_time, a.second_confirm_time, a.title, u.name as user_name").Join("inner", "users as u", "a.user_id=u.id").ID(id).Get(&auditInfo)
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	c.Data["auditInfo"] = auditInfo
	c.TplName = "audit/audit_detail.html"
}

// 全部列表
func (c *AuditController) AuditList() {
	page, _ := c.GetInt("page", 1)
	pagesize, _ := c.GetInt("pagesize", 10)

	offset := (page - 1) * pagesize

	var auditList []*AuditInfo
	err := db.Table("audit").Alias("a").Select("a.id, a.first_confirm, a.second_confirm, a.create_time, a.money, a.title, u.name as user_name").Join("inner", "users as u", "a.user_id=u.id").OrderBy("id desc").Limit(pagesize, offset).Find(&auditList)
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	total, err2 := db.Count(new(models.Audit))
	if err2 != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}

	c.Data["total"] = total
	c.Data["page"] = page
	c.Data["auditList"] = auditList

	c.TplName = "audit/audit_list.html"
}
