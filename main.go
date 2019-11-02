package main

import (
	"CoolQ/config"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	// 设置AppID  开发者域名反写.应用英文名
	cqp.AppID = "com.ypdan.ypdan"

	// 插件启动时被调用
	cqp.Enable = onEnable

	// 插件禁用时被调用
	cqp.Disable = onDisable

	// 注册接收到qq群消息的事件
	cqp.GroupMsg = onGroupMsg
}

//go:generate cqcfg .
// cqp: 名称: 优品单
// cqp: 版本: 1.1.0:1
// cqp: 作者: 孙龙
// cqp: 简介: 监听QQ群消息并POST到指定接口
func main() {}

var (
	conf = config.Config{}
)

func onDisable() int32 {
	printInfo("插件被禁用")
	return 0
}

// 当插件启用时被调用
func onEnable() int32 {
	defer handleErr()
	// 配置文件初始化
	c, err := config.LoadConfig()
	checkErr(-2, err)
	conf = *c

	printInfo("插件被启用")
	return 0
}

type PushGroupMessage struct {
	Body          string `json:"body"`
	QqGroupNumber int64  `json:"qqGroupNumber"`
	SendQQ        int64  `json:"sendQQ"`
	Timestamp     int64  `json:"timestamp"`
}

// 发送消息到 ypdan 服务
func sendMsg(pushGroupMessage PushGroupMessage) {
	pushGroupMessageJson, err := json.Marshal(pushGroupMessage)
	checkErr(-1, err)
	request, err := http.NewRequest(http.MethodPost, conf.MessageHandlerUrl, bytes.NewReader(pushGroupMessageJson))
	checkErr(1, err)
	defer request.Body.Close()

	request.Header.Set("signature", signData(conf.Token, &pushGroupMessage))
	request.Header.Set("content-type", "application/json;charset=UTF-8")
	response, err := http.DefaultClient.Do(request)
	checkErr(2, err)
	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}
}

func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	pushGroupMessage := PushGroupMessage{
		Body:          msg,
		QqGroupNumber: fromGroup,
		SendQQ:        fromQQ,
		Timestamp:     time.Now().UnixNano() / 1e6,
	}
	go sendMsg(pushGroupMessage)
	return 0
}

type byte2DSlice [][]byte

func (p byte2DSlice) Len() int {
	return len(p)
}
func (p byte2DSlice) Less(i, j int) bool {
	return bytes.Compare(p[i], p[j]) == -1
}
func (p byte2DSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

//生成签名
func gen(secret string, data ...[]byte) string {
	sort.Sort(byte2DSlice(data))
	buffer := bytes.Buffer{}
	buffer.Write([]byte(secret))
	for _, v := range data {
		buffer.Write(v)
	}
	h := md5.New()
	h.Write(buffer.Bytes())
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

func signData(signToken string, pushGroupMessage *PushGroupMessage) string {
	return gen(signToken,
		[]byte("body:"+pushGroupMessage.Body+":"),
		[]byte("qqGroupNumber:"+strconv.FormatInt(pushGroupMessage.QqGroupNumber, 10)+":"),
		[]byte("sendQQ:"+strconv.FormatInt(pushGroupMessage.SendQQ, 10)+":"),
		[]byte("timestamp:"+strconv.FormatInt(pushGroupMessage.Timestamp, 10)),
	)
}

// 抛异常
func checkErr(code int, err error) {
	if err != nil {
		printErr(code, err)
	}
}

func printErr(code int, err error) {
	cqp.AddLog(cqp.Error, fmt.Sprintf("错误 code:[%d]", code), err.Error())
}

func printInfo(msg string) {
	cqp.AddLog(cqp.Info, "通知消息", msg)
}

func handleErr() {
	if err := recover(); err != nil {
		cqp.AddLog(cqp.Fatal, "严重错误", fmt.Sprint(err))
	}
}
