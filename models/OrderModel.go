package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"sync"
)

//用户表
type Order struct {
	mu          sync.Mutex
	Id          int    `orm:"column(Id);unique"`
	Hash        string `orm:"unique;"`
	Value       string `orm:""`
	Gas         uint64 `orm:""`
	GasPrice    uint64 `orm:""`
	From        string `orm:""`
	BlockNumber int64  `orm:""`
	Nonce       uint64 `orm:""`
	Data        string `orm:""`
	To          string `orm:""`
	Timestamp   uint64 `orm:""`
	Status      string `orm:""`
}

func NewOrder() *Order {
	return &Order{}
}
func GetTableOrder() string {
	return getTable("order")
}

var (
	// 与变量对应的使用互斥锁
	countGuard sync.Mutex
)

type SyncMutex struct {
	or         *Order
	countGuard sync.Mutex
}

//批量插入
//@param            uid         interface{}         用户UID
//@return           UserInfo                        用户信息
func (u *Order) InsertAll(orders []*Order) bool {
	//var orders []*Order
	u.mu.Lock()
	qs := orm.NewOrm().QueryTable(GetTableOrder())
	i, _ := qs.PrepareInsert()
	for _, user := range orders {
		_, err := i.Insert(user)
		if err == nil {
			//fmt.Println("yi")
		}
	}
	i.Close() // 别忘记关闭 statement
	u.mu.Unlock()
	return true
}

// 单条插入
func (u *Order) Insert(order *Order) bool {
	u.mu.Lock()
	o := orm.NewOrm()
	_, err := o.Insert(order)
	if err == nil {
		return true
	} else {
		fmt.Println(err)
	}
	u.mu.Unlock()
	return false
}

// 删除
func (u *Order) Delete(number int64) bool {
	var mu sync.Mutex
	mu.Lock()
	//o := orm.NewOrm()
	qs := orm.NewOrm().QueryTable(GetTableOrder())
	if _, err := qs.Filter("block_number", number).Delete(); err == nil {
		return true
	} else {
		fmt.Println(err)
	}
	mu.Unlock()
	return false
}
