package main

import (
	"github.com/wq1019/ding_talk"
	"log"
)

const (
	dingApiUrl = "http://dev.hn-zm.com:8081/health_check"
)

var (
	dispatcher *Dispatcher

	MaxWorker = 2
	MaxQueue  = 10
)

func main() {
	for i := 0; i < 10; i++ {
		msg := ding_talk.MarkdownMessage{
			MsgType: ding_talk.Markdown,
			Markdown: ding_talk.MarkdownData{
				Title: "酷Q监控通知",
				Text:  "hello world",
			},
			At: &ding_talk.At{
				IsAtAll: true,
			},
		}
		work := Job{Payload: msg}
		JobQueue <- work
	}
	// 程序在这里等待
	for {
	}
}

func init() {
	// 创建调度器
	dispatcher = NewDispatcher(MaxWorker)
	// 调度器中创建n个work并且循环等待任务
	dispatcher.Run()

	// 创建任务队列
	JobQueue = make(chan Job, MaxQueue)

	Client = ding_talk.NewClient(dingApiUrl)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
