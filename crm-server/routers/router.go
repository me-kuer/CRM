package routers

import (
	"crm-server/controllers"
	"crm-server/models"
	"crm-server/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/url"
	"strconv"
)

func init() {
	//beego.Router("/", &controllers.MainController{})

	// 默认 / 重定向到 /index/index
	beego.Get("/", func(ctx *context.Context) {
		ctx.Redirect(302, "/index/index")
	})

	var indexNameSpace = beego.NewNamespace("/index",
		beego.NSBefore(checkLogin),
		beego.NSRouter("/index", &controllers.IndexController{}, "get:Index"),
		beego.NSRouter("/header", &controllers.IndexController{}, "get:Header"),
		beego.NSRouter("/main", &controllers.IndexController{}, "get:MainArea"),
		beego.NSRouter("/menu", &controllers.IndexController{}, "get:Menu"),
		beego.NSRouter("/welcome", &controllers.IndexController{}, "get:Welcome"),
	)

	beego.Router("/login", &controllers.UserController{}, "get:LoginPage;post:Login")
	beego.Router("/logout", &controllers.UserController{}, "get:Logout")

	var userNameSpace = beego.NewNamespace("/user",
		beego.NSBefore(checkLogin, checkRule),
		beego.NSRouter("/userlist", &controllers.UserController{}, "get:UserList"),
		beego.NSRouter("/add", &controllers.UserController{}, "get:AddPage;post:Add"),
		beego.NSRouter("/detail", &controllers.UserController{}, "get:Detail"),
		beego.NSRouter("/mod", &controllers.UserController{}, "post:Mod"),
		beego.NSRouter("/del", &controllers.UserController{}, "get:Del"),
	)

	var auditNameSpace = beego.NewNamespace("/audit",
		beego.NSBefore(checkLogin, checkRule), //过滤器
		beego.NSRouter("/add", &controllers.AuditController{}, "get:AddPage;post:Add"),
		beego.NSRouter("/detail", &controllers.AuditController{}, "get:Detail"),
		beego.NSRouter("/firstconfirmlist", &controllers.AuditController{}, "get:FirstConfirmList"),
		beego.NSRouter("/secondconfirmlist", &controllers.AuditController{}, "get:SecondConfirmList"),
		beego.NSRouter("/firstconfirm", &controllers.AuditController{}, "post:FirstConfirm"),
		beego.NSRouter("/secondconfirm", &controllers.AuditController{}, "post:SecondConfirm"),
		beego.NSRouter("/auditlist", &controllers.AuditController{}, "get:AuditList"),
	)

	var customerNameSpace = beego.NewNamespace("/customer",
		beego.NSBefore(checkLogin, checkRule), //过滤器
		beego.NSRouter("/add", &controllers.CustomerContoller{}, "get:AddPage;post:Add"),
		beego.NSRouter("/my", &controllers.CustomerContoller{}, "get:MyCustomer"),
		beego.NSRouter("/my_customer_detail", &controllers.CustomerContoller{}, "get:MyCustomerDetail"),
		beego.NSRouter("/mod_my_customer", &controllers.CustomerContoller{}, "post:ModMyCustomer"),
		beego.NSRouter("/public", &controllers.CustomerContoller{}, "get:PublicCustomer"),
		beego.NSRouter("/joinme", &controllers.CustomerContoller{}, "post:JoinMe"),
		beego.NSRouter("/customerlist", &controllers.CustomerContoller{}, "get:CustomerList"),
		beego.NSRouter("/detail", &controllers.CustomerContoller{}, "get:Detail"),
		beego.NSRouter("/mod", &controllers.CustomerContoller{}, "post:Mod"),
		beego.NSRouter("/del", &controllers.CustomerContoller{}, "get:Del"),
		beego.NSRouter("/excel_input", &controllers.CustomerContoller{}, "get:ExcelInput"),
	)
	// 错误页面
	beego.Router("/error", &controllers.InfoController{}, "get:ErrorPage")

	// 角色管理
	var roleNameSpace = beego.NewNamespace("/role",
		beego.NSBefore(checkLogin, checkRule), //过滤器
		beego.NSRouter("/rolelist", &controllers.RoleController{}, "get:RoleList"),
		beego.NSRouter("/add", &controllers.RoleController{}, "get:AddPage;post:Add"),
		beego.NSRouter("/del", &controllers.RoleController{}, "get:Del"),
		beego.NSRouter("/detail", &controllers.RoleController{}, "get:Detail"),
		beego.NSRouter("/mod", &controllers.RoleController{}, "post:Mod"),
	)
	// 注册路由分组
	beego.AddNamespace(indexNameSpace,userNameSpace,auditNameSpace,roleNameSpace,customerNameSpace)

}

// 验证用户是否登录
func checkLogin(ctx *context.Context) {
	// 如果没有获取到admin_id, 则跳转到登录页面
	_, ok := ctx.Input.Session("user_id").(int)
	if !ok {
		if ctx.Input.IsAjax() {
			data := map[string]string{
				"code": "401",
				"msg":  "对不起，您还没有登录",
			}
			ctx.Output.JSON(data, false, true)
			return
		} else {
			ctx.Redirect(302, "/login")
		}
	}
}

// 验证用户是否拥有该权限
func checkRule(ctx *context.Context) {
	userId := ctx.Input.Session("user_id").(int)

	// 如果用户为超级管理员(user_id=1)，则不进行验证权限
	if userId != 1 {
		// 查询该userId对应的角色下的角色id
		var user models.Users
		_, err := models.DB.Cols("role_id").ID(userId).Get(&user)
		if err != nil {
			utils.Logger.Error(err.Error())
			ctx.Redirect(302, "/error?msg="+err.Error())
			return
		}

		// 使用casbin验证
		// 获取请求方式
		//var method= ctx.Request.Method
		// 获取路由
		var uri= ctx.Request.RequestURI

		ok := models.Enforcer.Enforce(strconv.Itoa(user.RoleId), uri, "*")
		if !ok {
			if ctx.Input.IsAjax() {
				data := map[string]string{
					"code": "403",
					"msg":  "对不起，您的权限不足",
				}
				ctx.Output.JSON(data, false, true)
				return
			} else {
				ctx.Redirect(302, "/error?msg="+url.QueryEscape("对不起，您的权限不足"))
			}
		}
	}
}
