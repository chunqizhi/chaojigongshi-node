package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"os"
	"strings"
)

func init() {
	Init()
}

//初始化数据库注册
func Init() {
	//初始化数据库
	RegisterDB()
	runmode := beego.AppConfig.String("runmode")
	if runmode == "prod" {
		orm.Debug = false
		orm.RunSyncdb("default", false, false)
	} else {
		orm.Debug = true
		orm.RunSyncdb("default", false, true)
	}
}

//注册数据库
func RegisterDB() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	models := []interface{}{
		NewOrder(),
	}
	orm.RegisterModelWithPrefix(beego.AppConfig.DefaultString("prefix", "bee_"), models...)
	dbUser := beego.AppConfig.String("user")
	dbPassword := beego.AppConfig.String("password")
	if envpass := os.Getenv("MYSQL_PASSWORD"); envpass != "" {
		dbPassword = envpass
	}
	dbDatabase := beego.AppConfig.String("database")
	if envdatabase := os.Getenv("MYSQL_DATABASE"); envdatabase != "" {
		dbDatabase = envdatabase
	}
	dbCharset := beego.AppConfig.String("charset")
	dbHost := beego.AppConfig.String("host")
	if envhost := os.Getenv("MYSQL_HOST"); envhost != "" {
		dbHost = envhost
	}
	dbPort := beego.AppConfig.String("port")
	if envport := os.Getenv("MYSQL_PORT"); envport != "" {
		dbPort = envport
	}
	loc := "Local"
	if timezone := beego.AppConfig.String("timezone"); timezone != "" {
		loc = url.QueryEscape(timezone)
	}
	dbLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=%v", dbUser, dbPassword, dbHost, dbPort, dbDatabase, dbCharset, loc)
	maxIdle := beego.AppConfig.DefaultInt("maxIdle", 50)
	maxConn := beego.AppConfig.DefaultInt("maxConn", 300)
	fmt.Println("before:")
	fmt.Println("db link:",dbLink)
	fmt.Println("maxIdle:",maxIdle)
	fmt.Println("maxConn:",maxConn)
	dbLink = "root:root@tcp(myMysqlNode:3306)/etch?charset=utf8&loc=Local"
	maxIdle = 50
	maxConn = 300
	fmt.Println("db link:",dbLink)
	fmt.Println("maxIdle:",maxIdle)
	fmt.Println("maxConn:",maxConn)
	if err := orm.RegisterDataBase("default", "mysql", dbLink, maxIdle, maxConn); err != nil {
		panic(err)
	}
	db, _ := orm.GetDB("default")
	//设置连接池超时时间 mysql默认超时时间为28800秒也就是八个小时
	db.SetConnMaxLifetime(14400)
}

//获取带表前缀的数据表
//@param            table               数据表
func getTable(table string) string {
	prefix := beego.AppConfig.DefaultString("prefix", "bee_")
	if !strings.HasPrefix(table, prefix) {
		table = prefix + table
	}
	return table
}
