package main

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/chunqizhi/chaojigongshi-node/etch"
	"github.com/chunqizhi/chaojigongshi-node/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/garyburd/redigo/redis"
	"os"
	"time"
)

func startDO(client *etch.Eclient) bool {
	count, err := client.Count()
	if err != nil { //统计区块高度错误
		fmt.Println("client.Count(): ",err.Error())
		panic(err)
	}
	conn, err := redis.Dial("tcp", beego.AppConfig.String("redis"), redis.DialPassword(beego.AppConfig.String("redisPass")))
	if err != nil {
		fmt.Println("redis.Dial(): ",err.Error())
		panic(err)
		return false
	}
	num, err := redis.Int64(conn.Do("GET", "ethnum")) // 获取key
	if err != nil {                                   //没有这个key
		conn.Do("set", "ethnum", "0")
	}
	if num == 0 {
		num = 1
	}
	if num > count {
		fmt.Printf("num = %d, count = %d\n ",num,count)
		conn.Close()
		return false
	}
	block, err := client.Block(num)
	if err != nil { //查询区块详情错误
		fmt.Println("client.Block(): ",err.Error())
		fmt.Println("num = ",num)
		panic(err)
	}
	blocklen := len(block.Transactions())
	if blocklen <= 0 {
		fmt.Println("区块：", num, "没有交易数据")
		conn.Do("set", "ethnum", num+1) //退出之前区块高度加一
		conn.Close()
		return true
	}
	models.NewOrder().Delete(num) // 防止进程中断之后重启之后清理之前未完成的记录
	var orderPl []*models.Order
	var blockSum = (blocklen / 1000) * 1000 // 判断区块内数据是否满足一千条
	timestamp := uint64(block.Time())
	for key, tx := range block.Transactions() {
		receipt, err := client.GetTransactionReceipt(tx.Hash())
		if err != nil {
			fmt.Println("client.GetTransactionReceipt(): ",err.Error())
			if err == ethereum.NotFound {
				fmt.Printf("tx hash = %s\n",tx.Hash().String())
				continue
			}
			panic(err)
		}
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			fmt.Println("client.NetworkID(): ",err.Error())
			panic(err)
		}
		from := ""
		msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
		if err != nil {
			fmt.Println("tx.AsMessage(): ",err.Error())
			panic(err)
		}
		from = msg.From().String() // 获取from地址
		insert := new(models.Order)
		insert.Status = receipt.Status
		insert.Hash = tx.Hash().Hex()
		insert.Value = tx.Value().String()
		insert.Gas = tx.Gas()
		insert.From = from
		insert.BlockNumber = num
		insert.GasPrice = tx.GasPrice().Uint64()
		insert.Nonce = tx.Nonce()
		//insert.Data = tx.Data()
		if tx.To() != nil {
			insert.To = tx.To().String()
		}
		insert.Timestamp = timestamp
		if key+1%1000 == 0 {
			addSta := models.NewOrder().InsertAll(orderPl)
			if !addSta {
				fmt.Println("意外错误")
			} else {
				orderPl = orderPl[0:0] //清空切片
			}
		} else {
			if key+1 > blockSum { // 不足一千条的时，一条一条插入
				addSta := models.NewOrder().Insert(insert)
				if !addSta {
					fmt.Println(addSta)
					os.Exit(0)
				} else {
					orderPl = orderPl[0:0] //清空切片
				}
			} else {
				orderPl = append(orderPl, insert) //合并要插入数据库的数据
			}
		}
	}
	conn.Do("set", "ethnum", num+1) //退出之前区块高度加一
	conn.Close()
	fmt.Println("区块：", num, "记录完毕") // 1
	return true
}

// 自旋方法
func start(this *etch.Eclient) {
	status := startDO(this)
	if status {
		time.Sleep(5 * time.Millisecond)
		start(this)
	} else {
		fmt.Println("已扫描完所有区块，休眠一分钟,后继续执行")
		time.Sleep(1 * time.Minute)
		start(this)
	}
}
func main() {
	client, err := etch.New(beego.AppConfig.String("ethNode")) // 只在开始时候开启一次链接，
	if err != nil {                                            //捕获到异常，可能是端口耗尽，休眠一段实践进行重试
		fmt.Println("休眠30秒后重启主进程")
		time.Sleep(30 * time.Second) // tcp端口默认失效时间最低30秒
		main()
	}
	defer func() {
		if r := recover(); r != nil { // 捕获到异常
			client.Close() // 关闭链接
			main()         //重启主方法
		}
	}()
	start(client)
}
