package controllers

import (
	"crm-server/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
)

type CustomerContoller struct {
	beego.Controller
}

type CustomerInfo struct {
	models.Customer `xorm:"extends"`
	UserName        string `xorm:""`
}

// 添加客户
func (c *CustomerContoller) Add() {
	customerListJson := c.GetString("customer_list", "[]")

	// 将json转为结构体
	var customerListMap []map[string]interface{}
	err := json.Unmarshal([]byte(customerListJson), &customerListMap)
	fmt.Println(customerListMap)
	if err != nil {
		log.Error(err.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err.Error(),
		}
		c.ServeJSON()
		return
	}
	// map 转为 struct
	var customerList []models.Customer
	for _, v := range customerListMap {
		name := v["name"].(string)
		company := v["company"].(string)
		tel := v["tel"].(string)
		addr := v["addr"].(string)
		tag, _ := strconv.Atoi(v["tag"].(string))
		stage, _ := strconv.Atoi(v["stage"].(string))
		contract, _ := strconv.Atoi(v["contract"].(string))
		price, _ := strconv.ParseFloat(v["price"].(string), 64)
		amount, _ := strconv.ParseFloat(v["amount"].(string), 64)
		appointStart := v["appoint_start"].(string)
		appointEnd := v["appoint_end"].(string)
		period := v["period"].(string)
		payee := v["payee"].(string)
		payeeUsername := v["payee_username"].(string)
		bank := v["bank"].(string)
		remark := v["remark"].(string)

		var customer = models.Customer{
			Name:          name,
			Company:       company,
			Tel:           tel,
			Addr:          addr,
			Tag:           tag,
			Stage:         stage,
			Contract:      contract,
			Price:         price,
			Amount:        amount,
			AppointStart:  appointStart,
			AppointEnd:    appointEnd,
			Period:        period,
			Payee:         payee,
			PayeeUsername: payeeUsername,
			Bank:          bank,
			Remark:        remark,
			CreateTime:    strconv.FormatInt(time.Now().Unix(), 10),
			Status:        1,
		}
		customerList = append(customerList, customer)
	}

	_, err2 := db.Insert(&customerList)
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

// 添加客户页面
func (c *CustomerContoller) AddPage() {
	c.TplName = "customer/add_customer.html"
}

// 加入到我的客户列表
func (c *CustomerContoller) JoinMe() {
	customerIdList := c.GetString("customer_list", "")
	// session中获取我的user_id
	userId := c.GetSession("user_id")
	var customer = models.Customer{
		UserId: userId.(int),
	}

	row, err := db.Where("user_id=?", 0).In("id", strings.Split(customerIdList, ",")).Update(&customer)

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
			"msg":  "添加失败",
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

// 移除客户
func (c *CustomerContoller) Del() {
	// 将客户的status 变为已删除
	id, _ := c.GetInt("id")
	var customer = models.Customer{
		Status: 0,
	}

	row, err := db.Id(id).Cols("status").Update(&customer)
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

// 我的客户详情
func (c *CustomerContoller) MyCustomerDetail() {
	id, err := c.GetInt("id")
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	var customer models.Customer
	userId := c.GetSession("user_id")

	_, err2 := db.Id(id).Where("user_id=?", userId).Get(&customer)
	if err2 != nil {
		log.Error(err2.Error())
		c.Redirect("/error?msg="+err2.Error(), 302)
		return
	}

	c.Data["customer"] = customer

	c.TplName = "customer/my_customer_detail.html"
}

// 修改我的客户
func (c *CustomerContoller) ModMyCustomer() {
	id, err := c.GetInt("id")
	if err != nil {
		log.Error(err.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err.Error(),
		}
		c.ServeJSON()
		return
	}
	name := c.GetString("name", "")
	company := c.GetString("company", "")
	tel := c.GetString("tel", "")
	addr := c.GetString("addr", "")
	tag, _ := c.GetInt("tag", 1)
	stage, _ := c.GetInt("stage", 1)
	contract, _ := c.GetInt("contract", 1)
	price, _ := c.GetFloat("price", 0)
	amount, _ := c.GetFloat("amount", 0)
	appointStart := c.GetString("appoint_start", "0")
	appointEnd := c.GetString("appoint_end", "0")
	period := c.GetString("period", "")
	payee := c.GetString("payee", "")
	payeeUsername := c.GetString("payee_username", "")
	bank := c.GetString("bank", "")
	remark := c.GetString("remark", "")

	// 每修改一次需要更新一次
	var customer = models.Customer{
		Name:          name,
		Company:       company,
		Tel:           tel,
		Addr:          addr,
		Tag:           tag,
		Stage:         stage,
		Contract:      contract,
		Price:         price,
		Amount:        amount,
		AppointStart:  appointStart,
		AppointEnd:    appointEnd,
		Period:        period,
		Payee:         payee,
		PayeeUsername: payeeUsername,
		Bank:          bank,
		Remark:        remark,
		UpdateTime:    strconv.FormatInt(time.Now().Unix(), 10),
	}
	// 只能修改我自己名下的客户信息
	userId := c.GetSession("user_id")
	row, err2 := db.Id(id).Where("user_id=?", userId).Update(&customer)
	if err2 != nil {
		log.Error(err2.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err2.Error(),
		}
		c.ServeJSON()
		return
	}

	if row < 1 {
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "修改失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 200,
		"msg":  "修改成功",
	}
	c.ServeJSON()
}

// 所有客户列表
func (c *CustomerContoller) CustomerList() {
	page, _ := c.GetInt("page", 1)
	pagesize, _ := c.GetInt("pagesize", 10)
	offset := (page - 1) * pagesize

	var customerList []*CustomerInfo
	/*
		select
			c.id, c.name, c.stage, c.addr, c.tag, c.company, c.tel, c.create_time, u.name as user_name
		from
			customer as c
		left join
			users as u
		on
			c.user_id = u.id
		where
			c.status = 1
		order by
			c.id desc
		limit 1,10

	*/
	sql := "select c.id, c.name, c.addr, c.stage, c.tag, c.company, c.tel, c.create_time, u.name as user_name from customer as c left join users as u on c.user_id = u.id where c.status=1 order by c.id desc limit ?,?"
	err := db.SQL(sql, offset, pagesize).Find(&customerList)

	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}

	var customer models.Customer
	total, err2 := db.Where("status=1").Count(&customer)
	if err2 != nil {
		log.Error(err2.Error())
		c.Redirect("/error?msg="+err2.Error(), 302)
		return
	}
	c.Data["total"] = total
	c.Data["customerList"] = customerList

	c.TplName = "customer/customer_list.html"
}

// 公共客户列表
func (c *CustomerContoller) PublicCustomer() {
	page, _ := c.GetInt("page", 1)
	pagesize, _ := c.GetInt("pagesize", 10)
	offset := (page - 1) * pagesize

	var customerList []*models.Customer
	err := db.Cols("id", "name", "stage", "company", "tel", "create_time").Where("user_id=0 and status=1").Limit(pagesize, offset).Find(&customerList)
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	var customer models.Customer
	total, err2 := db.Where("user_id=0 and status=1").Count(&customer)
	if err2 != nil {
		log.Error(err2.Error())
		c.Redirect("/error?msg="+err2.Error(), 302)
		return
	}
	c.Data["total"] = total
	c.Data["customerList"] = customerList

	c.TplName = "customer/public_customer.html"
}

// 我的客户列表
func (c *CustomerContoller) MyCustomer() {
	keywords := c.GetString("keywords", "")
	stage := c.GetString("stage", "")
	tag := c.GetString("tag", "")
	appointEnd := c.GetString("appointtime", "0")
	updateTime := c.GetString("updatetime", "0")
	page, _ := c.GetInt("page", 1)
	pagesize, _ := c.GetInt("pagesize", 10)
	offset := (page - 1) * pagesize

	userId := c.GetSession("user_id")
	where := "user_id=" + strconv.Itoa(userId.(int)) + " and appoint_end>=" + appointEnd + " and update_time>=" + updateTime + " and status=1 "

	if keywords != "" {
		where += " and (name like '%" + keywords + "%' or addr like '%" + keywords + "%')"
	}
	if tag != "" {
		where += " and tag=" + tag
	}
	if stage != "" {
		where += " and stage=" + stage
	}

	var customerList []*models.Customer
	err := db.Cols("id", "tel", "name", "company", "addr", "create_time", "appoint_end", "tag", "stage", "update_time").Where(where).OrderBy("id desc").Limit(pagesize, offset).Find(&customerList)
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}

	total, err2 := db.Where(where).Count(new(models.Customer))
	if err2 != nil {
		log.Error(err2.Error())
		c.Redirect("/error?msg="+err2.Error(), 302)
		return
	}

	c.Data["customerList"] = customerList
	c.Data["total"] = total

	c.TplName = "customer/my_customer.html"

}

// 客户详情
func (c *CustomerContoller) Detail() {
	id, err := c.GetInt("id")
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	var customer models.Customer

	_, err2 := db.Id(id).Get(&customer)
	if err2 != nil {
		log.Error(err2.Error())
		c.Redirect("/error?msg="+err2.Error(), 302)
		return
	}

	c.Data["customer"] = customer

	c.TplName = "customer/customer_detail.html"
}

// 修改客户
func (c *CustomerContoller) Mod() {
	id, err := c.GetInt("id")
	if err != nil {
		log.Error(err.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err.Error(),
		}
		c.ServeJSON()
		return
	}
	name := c.GetString("name", "")
	company := c.GetString("company", "")
	tel := c.GetString("tel", "")
	addr := c.GetString("addr", "")
	tag, _ := c.GetInt("tag", 1)
	stage, _ := c.GetInt("stage", 1)
	contract, _ := c.GetInt("contract", 1)
	price, _ := c.GetFloat("price", 0)
	amount, _ := c.GetFloat("amount", 0)
	appointStart := c.GetString("appoint_start", "0")
	appointEnd := c.GetString("appoint_end", "0")
	period := c.GetString("period", "")
	payee := c.GetString("payee", "")
	payeeUsername := c.GetString("payee_username", "")
	bank := c.GetString("bank", "")
	remark := c.GetString("remark", "")

	// 每修改一次需要更新一次
	var customer = models.Customer{
		Name:          name,
		Company:       company,
		Tel:           tel,
		Addr:          addr,
		Tag:           tag,
		Stage:         stage,
		Contract:      contract,
		Price:         price,
		Amount:        amount,
		AppointStart:  appointStart,
		AppointEnd:    appointEnd,
		Period:        period,
		Payee:         payee,
		PayeeUsername: payeeUsername,
		Bank:          bank,
		Remark:        remark,
		UpdateTime:    strconv.FormatInt(time.Now().Unix(), 10),
	}
	// 只能修改我自己名下的客户信息
	row, err2 := db.Id(id).Update(&customer)
	if err2 != nil {
		log.Error(err2.Error())
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  err2.Error(),
		}
		c.ServeJSON()
		return
	}

	if row < 1 {
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "修改失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 200,
		"msg":  "修改成功",
	}
	c.ServeJSON()
}


func (c *CustomerContoller) ExcelInput(){
	c.TplName = "customer/excel_input.html"
}