package controllers

import (
	"crm-server/models"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

type UserInfo struct {
	models.Users `xorm:"extends"`
	RoleName     string `xorm:""`
}

// 用户登录页面
func (c *UserController) LoginPage() {
	c.TplName = "user/login.html"
}

// 用户登录
func (c *UserController) Login() {
	// 获取用户名、密码
	username := c.GetString("username")
	password := c.GetString("password")

	// MD5加密password
	h := md5.New()
	h.Write([]byte(password))
	encPassword := hex.EncodeToString(h.Sum(nil))

	// 查询user_id
	var user models.Users
	_, err := db.Cols("id","status").Where("username=? and password=?", username, encPassword).Get(&user)
	if err != nil {
		// 返回错误
		res := map[string]string{
			"code": "500",
			"msg":  err.Error(),
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if user.Id <= 0 {
		c.Data["json"] = map[string]string{
			"code": "500",
			"msg":  "用户名或密码错误！",
		}
		c.ServeJSON()
		return
	}
	// 判断用户是否被禁止登陆
	if user.Status == 2 {
		c.Data["json"] = map[string]string{
			"code": "500",
			"msg":  "对不起，您已被禁止登陆",
		}
		c.ServeJSON()
		return
	}
	
	// 用户名密码正确，session保存user_id
	c.SetSession("user_id", user.Id)

	// 返回json
	c.Data["json"] = map[string]string{
		"code": "200",
		"msg":  "登录成功",
	}
	c.ServeJSON()
}

// 用户退出
func (c *UserController) Logout() {
	// 删除session
	c.DelSession("user_id")
}

// 用户列表
func (c *UserController) UserList() {
	var userList []*UserInfo
	// 用户表 & 角色表联查
	err := db.Table("users").Alias("u").Select("u.id,u.name, u.tel, u.username, r.name as role_name").Join("left", "role as r", "u.role_id=r.id").In("status", []int{1, 2}).Find(&userList)
	if err != nil {
		log.Error(err.Error())
		// 跳转到失败页
		c.Redirect("/error?msg="+err.Error(), 302)

		return
	}
	c.Data["userList"] = userList
	c.TplName = "user/user_list.html"
}

// 用户详情
func (c *UserController) Detail() {
	user_id, err := c.GetInt("user_id")
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}

	// 获取用户详情信息
	var user models.Users
	_, err2 := db.Cols("id", "name", "username", "tel", "role_id", "status").ID(user_id).Get(&user)

	if err2 != nil {
		log.Error(err2.Error())
		c.Redirect("/error?msg="+err2.Error(), 302)
		return
	}
	// 获取角色列表
	var roleList []*models.Role
	err3 := db.Cols("id", "name").Find(&roleList)
	if err3 != nil {
		log.Error(err3.Error())
		c.Redirect("/error?msg="+err3.Error(), 302)
		return
	}
	// 输出模板
	c.Data["user"] = user
	c.Data["roleList"] = roleList
	c.TplName = "user/user_detail.html"
}

// 修改用户
func (c *UserController) Mod() {
	id, _ := c.GetInt("id")
	name := c.GetString("name")
	tel := c.GetString("tel")
	username := c.GetString("username")
	password := c.GetString("password")
	status, _ := c.GetInt("status", 1)
	fmt.Println(status)
	role, _ := c.GetInt("role", 0)

	var user = models.Users{
		Name:     name,
		Tel:      tel,
		Username: username,
		Status:   status,
		RoleId:   role,
	}
	if password != "" {
		// MD5加密password
		h := md5.New()
		h.Write([]byte(password))
		encPassword := hex.EncodeToString(h.Sum(nil))

		user.Password = encPassword
	}
	_, err := db.Id(id).Update(&user)
	if err != nil {
		log.Error(err.Error())
		// 返回json
		c.Data["json"] = map[string]string{
			"code": "500",
			"msg":  err.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{
		"code": "200",
		"msg":  "保存成功",
	}
	c.ServeJSON()
}

// 添加用户页面
func (c *UserController) AddPage() {
	// 获取角色列表
	var roleList []*models.Role
	err := db.Cols("id", "name").Find(&roleList)
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	c.Data["roleList"] = roleList
	c.TplName = "user/add_user.html"
}

// 添加用户提交
func (c *UserController) Add() {
	name := c.GetString("name")
	tel := c.GetString("tel")
	username := c.GetString("username")
	password := c.GetString("password")
	status, _ := c.GetInt("status", 1)
	role, _ := c.GetInt("role", 0)

	// MD5加密password
	h := md5.New()
	h.Write([]byte(password))
	encPassword := hex.EncodeToString(h.Sum(nil))

	var user = models.Users{
		Name:     name,
		Tel:      tel,
		Username: username,
		Password: encPassword,
		Status:   status,
		RoleId:   role,
	}

	_, err := db.Insert(&user)
	if err != nil {
		log.Error(err.Error())
		// 返回json
		c.Data["json"] = map[string]string{
			"code": "500",
			"msg":  err.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{
		"code": "200",
		"msg":  "添加用户成功",
	}
	c.ServeJSON()
}

// 删除用户
func (c *UserController) Del() {
	user_id, err := c.GetInt("id")
	if err != nil {
		log.Error(err.Error())
		c.Redirect("/error?msg="+err.Error(), 302)
		return
	}
	//创建事务会话
	dbSession := db.NewSession()
	// 预关闭事务
	defer dbSession.Close()
	// 开启事务
	err2 := dbSession.Begin()
	if err2 != nil {
		msg := err2.Error()
		log.Error(msg)
		res := map[string]string{
			"code": "500",
			"msg":  msg,
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 将用户的status 改为0(删除状态)
	var user = models.Users{
		Status: 0,
	}
	rows, err2 := dbSession.ID(user_id).Cols("status").Update(&user)
	if err2 != nil {
		// 回滚
		dbSession.Rollback()

		log.Error(err2.Error())
		// 返回json
		c.Data["json"] = map[string]string{
			"code": "500",
			"msg":  err2.Error(),
		}
		c.ServeJSON()
		return
	}
	// 修改的记录数为0条，返回 删除失败
	if rows < 1 {
		// 回滚
		dbSession.Rollback()
		// 返回json
		c.Data["json"] = map[string]string{
			"code": "500",
			"msg":  "删除失败",
		}
		c.ServeJSON()
		return
	}

	// 将该用户名下的客户，未签约转移为公共客户下
	sql := "update customer set user_id=0 where user_id=? and stage!=3"
	_, err3 := dbSession.Exec(sql, user_id)
	if err3 != nil {
		// 回滚
		dbSession.Rollback()

		log.Error(err3.Error())
		// 返回json
		c.Data["json"] = map[string]string{
			"code": "500",
			"msg":  err3.Error(),
		}
		c.ServeJSON()
		return
	}
	// 已签约的转移到总管理名下
	sql2 := "update customer set user_id=1 where user_id=? and stage=3"
	_, err4 := dbSession.Exec(sql2, user_id)
	if err4 != nil {
		// 回滚
		dbSession.Rollback()

		log.Error(err4.Error())
		// 返回json
		c.Data["json"] = map[string]string{
			"code": "500",
			"msg":  err4.Error(),
		}
		c.ServeJSON()
		return
	}
	// 提交事务
	dbSession.Commit()

	c.Data["json"] = map[string]string{
		"code": "200",
		"msg":  "删除成功",
	}
	c.ServeJSON()
	return
}
