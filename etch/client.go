package etch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type result struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int64       `json:"id"`
	Result  interface{} `json:"result"`
}

var url = "http://localhost:18545"
var contentType = "application/json;charset=utf-8"

type Chao struct {
}

func (c *Chao) GetBlockNumber() int{
	param := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`
	re := c.Post(param)
	ree := re.Result.(string)
	i, err := strconv.ParseInt(ree[2:], 16, 64)
	if err != nil {
		fmt.Println("err: ", err.Error())
		return 0
	}
	return int(i)
}

func (c *Chao) GetBlockByNumber(blockNumber int) map[string]interface{} {

	bn := fmt.Sprintf("%x", blockNumber)

	bnh := "0x" + bn

	param := "{\"jsonrpc\":\"2.\",\"method\":\"eth_getBlockByNumber\",\"params\":[\"" + bnh + "\", true],\"id\":1}"

	re := c.Post(param)

	return re.Result.(map[string]interface{})
}


func (c *Chao) GetTransactionReceipt(hash string) {
	param := "{\"jsonrpc\":\"2.0\",\"method\":\"eth_GetTransactionReceipt\",\"params\":[\""+ hash +"\"],\"id\":1}"
	re := c.Post(param)
	fmt.Println("re==",re)
}

func (c *Chao) GetTransactionByHash(hash string) {
	param := "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionByHash\",\"params\":[\""+ hash +"\"],\"id\":1}"
	re := c.Post(param)
	fmt.Println("re==",re)
}


func (c *Chao) Post(param string) *result {
	fmt.Println("lalalalala6")

	resp, err := http.Post(url, contentType, strings.NewReader(param))
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("err: ", err.Error())
		return nil
	}
	fmt.Println("lalalalala7")

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err: ", err.Error())
		return nil
	}
	fmt.Println("lalalalala8")

	re := &result{}
	err = json.Unmarshal(bytes, re)
	if err != nil {
		fmt.Println("err: ", err.Error())
		return nil
	}
	return re
}
