package main

import (
	"CoolQ/config"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/wq1019/ding_talk"
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
	// 注册接收到qq群消息的事件
	cqp.GroupMsg = onGroupMsg
	// 插件禁用时被调用
	cqp.Disable = onDisable
	cqp.Start = onStart
	cqp.Exit = onExit

	c, err := config.LoadConfig()
	checkErr(err)
	conf = *c
	b := ding_talk.NewClient(conf.NotifyUrl)
	dingTalkClient = *b
}

//go:generate cqcfg .
// cqp: 名称: ypdan
// cqp: 版本: 1.0.0:1
// cqp: 作者: sunlong
// cqp: 简介: 监听QQ群消息并POST到指定接口
func main() {}

var (
	conf           = config.Config{}
	dingTalkClient = ding_talk.DingTalkClient{}
)

func onStart() int32 {
	notifyDingDing(0, 0, "酷Q正在启动", "")
	return 0
}

func onExit() int32 {
	notifyDingDing(0, 0, "酷Q正在关闭", "")
	return 0
}

// 当插件启用时被调用
func onEnable() int32 {
	defer handleErr()
	notifyDingDing(0, 0, "插件被启用", "")
	return 0
}

type PushGroupMessage struct {
	Body          string `json:"body"`
	QqGroupNumber int64  `json:"qqGroupNumber"`
	SendQQ        int64  `json:"sendQQ"`
	Timestamp     int64  `json:"timestamp"`
}

func onDisable() int32 {
	notifyDingDing(0, 0, "插件被禁用", "")
	return 0
}

func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	pushGroupMessage := &PushGroupMessage{
		Body:          msg,
		QqGroupNumber: fromGroup,
		SendQQ:        fromQQ,
		Timestamp:     time.Now().UnixNano() / 1e6,
	}
	pushGroupMessageJson, err := json.Marshal(pushGroupMessage)
	checkErr(err)

	request, err := http.NewRequest(http.MethodPost, conf.MessageHandlerUrl, bytes.NewReader(pushGroupMessageJson))
	checkErr(err)
	defer request.Body.Close()

	request.Header.Set("signature", signData(conf.Token, pushGroupMessage))
	request.Header.Set("content-type", "application/json;charset=UTF-8")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		printErr(err)
		notifyDingDing(fromGroup, fromQQ, fmt.Sprintf("推送商品消息到优品单服务器失败, err: %s", err.Error()), msg)
	} else {
		_ = response.Body.Close()
	}
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

func timeFormat(timeInt int64) string {
	t := time.Unix(timeInt, 0)
	return fmt.Sprintf("%d月%d日%d时%d分%d秒", t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func notifyDingDing(fromGroup, fromQQ int64, msg, data string) {
	cqp.AddLog(cqp.Debug, "test-test", "sssss")
	markdown := ding_talk.MarkdownMessage{
		MsgType: ding_talk.Markdown,
		Markdown: ding_talk.MarkdownData{
			Title: "酷Q监控通知",
			Text: fmt.Sprintf("#### 酷Q监控通知\n"+
				"> FromGroup: %d\n\n"+
				"> FromQQ: %d\n\n"+
				"> Message: %s\n\n"+
				"> PostData: %s\n\n"+
				"> ###### %s发布 [优品单](https://ypdan.com) \n", fromGroup, fromQQ, msg, data, timeFormat(time.Now().Unix())),
		},
		At: &ding_talk.At{
			IsAtAll: true,
		},
	}
	_, err := dingTalkClient.Execute(markdown)
	cqp.AddLog(cqp.Debug, "test-1", "sssss")
	checkErr(err)
	cqp.AddLog(cqp.Debug, "test-2", "sssss")
}

func checkErr(err error) {
	if err != nil {
		printErr(err)
	}
}

func printErr(err error) {
	cqp.AddLog(cqp.Error, "错误", err.Error())
}

func handleErr() {
	if err := recover(); err != nil {
		cqp.AddLog(cqp.Fatal, "严重错误", fmt.Sprint(err))
	}
}
