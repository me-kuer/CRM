# 用户表
create table users (
    id int auto_increment primary key,
    username varchar(20) not null default "" comment "用户名",
    password varchar(32) not null default "" comment "密码",
    name varchar(60) not null default "" comment "姓名",
    tel varchar(11) default "" comment "联系方式",
    role_id int default 0 comment "角色id",
    status tinyint(1) unsigned not null default 1 comment "状态(1允许登录，0已删除，2禁止登录)"
) comment "用户表";

# 角色表
create table role(
    id int auto_increment primary key,
    name varchar(30) not null default "" comment "角色名称",
    modules varchar(300) default "" comment "模块id列表",
    desc text comment "描述"
) comment "角色表";

# 权限表( 使用casbin自带 )

# 模块表
create table modules(
    id int auto_increment primary key,
    name varchar(60) not null default "" comment "模块名称"
) comment "权限模块表";


# 客户表
create table customer(
    id int auto_increment primary key,
    name varchar(60) not null default "" comment "姓名",
    company varchar(180) default "" comment "公司（或证书）",
    tel varchar(21) default "" comment "联系方式",
    addr varchar(120) default "" comment "地址",
    tag tinyint(1) unsigned not null default 1 comment "标签（1人才库，2企业库）",
    stage tinyint(1) unsigned not null default 1 comment "阶段（0失效，1入库，2沟通中，3签约）",
    contract tinyint(1) unsigned not null default 1 comment "合同类型（1纸质合同）",
    price float(10,2) unsigned default 0 comment "价格",
    amount float(10,2) unsigned default 0 comment "付款金额",
    appoint_start varchar(11) default "0" comment "签约开始日期",
    appoint_end varchar(11) default "0" comment "签约结束日期",
    period varchar(60) default "" comment "签约期限",
    payee varchar(60) default "" comment "收款人",
    payee_username varchar(60) default "" comment "收款账户",
    bank varchar(60) default "" comment "银行",
    remark text comment "备注",
    user_id int default 0 comment "用户id",
    create_time varchar(11) default "0" comment "创建时间",
    status tinyint(1) unsigned default 1 comment "状态（1正常，0已删除）"
) comment "客户表";

# 结款审批表
create table audit(
    id int auto_increment primary key,
    user_id int not null default 0 comment "用户id",
    title varchar(60) not null default "" comment "标题",
    money float(10,2) not null default 0 comment "金额",
    desc text comment "描述",
    first_confirm tinyint(1) not null default 2 comment "第一级确认（1通过，0未审核，2拒绝）",
    first_confirm_time varchar(11) not null default "0" comment "第一级确认时间",
    second_confirm tinyint(1) not null default 2 comment "第二级确认（1通过，0未审核，2拒绝）",
    second_confirm_time varchar(11) not null default "0" comment "第二级确认时间"
    create_time varchar(11) not null default "0" comment "发起审批时间"
) comment "结款审批表";