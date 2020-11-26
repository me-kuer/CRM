package models

import (
	"crm-server/utils"
	"github.com/astaxie/beego"
	"github.com/casbin/casbin"
	"github.com/casbin/xorm-adapter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"os"
	"xorm.io/core"
)

// 用户表
type Users struct {
	Id       int    `xorm:"pk autoincr" json:"id,omitempty"`
	Username string `xorm:"varchar(20) notnull default('')" json:"username,omitempty"`
	Password string `xorm:"varchar(32) notnull default('')" json:"password,omitempty"`
	Name     string `xorm:"varchar(60) notnull default('')" json:"name,omitempty"`
	Tel      string `xorm:"varchar(11) default('')" json:"tel,omitempty"`
	RoleId   int    `xorm:"default(0)" json:"role_id,omitempty"`
	Status   int    `xorm:"tinyint(1) notnull default(1)" json:"status,omitempty"`
}

// 角色表
type Role struct {
	Id      int    `xorm:"pk autoincr" json:"id,omitempty"`
	Name    string `xorm:"varchar(30) notnull default('')" json:"name,omitempty"`
	Modules string `xorm:"varchar(300) default('0')" json:"modules,omitempty"`
	Desc    string `xorm:"text" json:"desc,omitempty"`
}

// 客户表
type Customer struct {
	Id            int     `xorm:"pk autoincr" json:"id,omitempty"`
	Name          string  `xorm:"varchar(60) notnull default('')" json:"name,omitempty"`
	Company       string  `xorm:"varchar(180) default('')" json:"company,omitempty"`
	Tel           string  `xorm:"varchar(21) default('')" json:"tel,omitempty"`
	Addr          string  `xorm:"varchar(120) default('')" json:"addr,omitempty"`
	Tag           int     `xorm:"tinyint(1) notnull default(1)" json:"tag,omitempty"`
	Stage         int     `xorm:"tinyint(1) notnull default(1)" json:"stage,omitempty"`
	Contract      int     `xorm:"tinyint(1) unsigned not null default(1)" json:"contract,omitempty"`
	Price         float64 `xorm:"float(10,2) default(0)" json:"price,omitempty"`
	Amount        float64 `xorm:"float(10,2) default(0)" json:"amount,omitempty"`
	AppointStart  string  `xorm:"varchar(11) default('0')" json:"appoint_start,omitempty"`
	AppointEnd    string  `xorm:"varchar(11) default('0')" json:"appoint_end,omitempty"`
	Period        string  `xorm:"varchar(60) default('')" json:"period,omitempty"`
	Payee         string  `xorm:"varchar(60) default('')" json:"payee,omitempty"`
	PayeeUsername string  `xorm:"varchar(60) default('')" json:"payee_username,omitempty"`
	Bank          string  `xorm:"varchar(60) default('')" json:"bank,omitempty"`
	Remark        string  `xorm:"text" json:"remark,omitempty"`
	UserId        int     `xorm:"default(0)" json:"user_id,omitempty"`
	UpdateTime    string  `xorm:"varchar(11) notnull default('0')" json:"update_time,omitempty"`
	CreateTime    string  `xorm:"varchar(11) notnull default('0')" json:"create_time,omitempty"`
	Status        int     `xorm:"tinyint(1) notnull default(1)" json:"status,omitempty"`
}

// 结款审批表
type Audit struct {
	Id                int     `xorm:"pk autoincr" json:"id,omitempty"`
	UserId            int     `xorm:"notnull default(0)" json:"user_id,omitempty"`
	Title             string  `xorm:"varchar(60) notnull default('')" json:"title,omitempty"`
	Money             float64 `xorm:"float(10,2) notnull default(0)" json:"money,omitempty"`
	Desc              string  `xorm:"text"`
	FirstConfirm      int     `xorm:"tinyint(1) notnull default(0)" json:"first_confrim,omitempty"`
	FirstConfirmTime  string  `xorm:"varchar(11) notnull default('0')" json:"first_confirm_time,omitempty"`
	SecondConfirm     int     `xorm:"tinyint(1) notnull default(2)" json:"second_confirm,omitempty"`
	SecondConfirmTime string  `xorm:"varchar(11) notnull default('0')" json:"second_confirm_time,omitempty"`
	CreateTime        string  `xorm:"varchar(11) notnull default(0)" json:"create_time,omitempty"`
}

// 权限模块表
type Modules struct {
	Id   int    `xorm:"pk autoincr" json:"id,omitempty"`
	Name string `xorm:"varchar(30) default('')" json:"name,omitempty"`
}

var DB *xorm.Engine

var Enforcer *casbin.Enforcer

// 日志引擎
var log = utils.Logger

func init() {
	// 读取数据库配置
	var (
		dbDriver   = beego.AppConfig.String("db_driver")
		dbName     = beego.AppConfig.String("db_name")
		dbHost     = beego.AppConfig.String("db_host")
		dbPort     = beego.AppConfig.String("db_port")
		dbUsername = beego.AppConfig.String("db_username")
		dbPassword = beego.AppConfig.String("db_password")
		dbPrefix   = beego.AppConfig.String("db_prefix")
	)
	var err error
	// 创建xorm引擎
	DB, err = xorm.NewEngine(dbDriver, dbUsername+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8")
	if err != nil {
		// 日志输出错误
		log.Error(err.Error())
		return
	}

	// 设置xorm日志
	f, err := os.Create("logs/xorm.log")
	if err != nil {
		println(err.Error())
		return
	}
	DB.SetLogger(xorm.NewSimpleLogger(f))

	// 延迟关闭不能用，否则会报 database is closed！
	//defer Engine.Close()

	err2 := DB.Ping()
	if err2 != nil {
		// 日志输出错误
		log.Error(err2.Error())
		return
	}
	// 打印sql语句
	DB.ShowSQL(true)

	// 设置日志等级
	DB.Logger().SetLevel(core.LOG_DEBUG)

	// 设置映射方式
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, dbPrefix)
	DB.SetTableMapper(tbMapper)

	// 创建表
	err3 := DB.Sync2(new(Users), new(Role), new(Customer), new(Modules), new(Audit))
	if err3 != nil {
		// 日志输出错误
		log.Error(err3.Error())
		return
	}

	// 初始化casbin
	a := xormadapter.NewAdapterByEngine(DB)

	Enforcer = casbin.NewEnforcer("conf/rbac_models.conf", a)

	//从mysql中加载策略
	err4 := Enforcer.LoadPolicy()
	if err4 != nil {
		log.Error(err4.Error())
		return
	}

	// 判断是否有超级管理员，如果没有则进行创建
	var user Users
	_,err5 := DB.Id(1).Get(&user)
	if err5 != nil {
		log.Error(err5.Error())
	}
	// 不存在进行创建
	if user.Id <= 0 {
		var admin = Users{
			Id: 1,
			Name: "超级管理员",
			Username: "admin",
			Password: "e10adc3949ba59abbe56e057f20f883e",
			Status: 1,
		}
		_,err6 := DB.Insert(&admin)
		if err6 != nil {
			log.Error(err6.Error())
		}
	}

}
