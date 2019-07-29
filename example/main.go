package main

import (
	"fmt"
	"github.com/wq1019/ding_talk"
	"time"
)

var ch2 = make(chan string, 2)

func main() {

	//ch := make(chan int, 3)
	//ch1 := make(chan string, 3)

	count := 6

	go func3()

	for i := 0; i < count; i++ {
		ch2 <- fmt.Sprintf("hello - %d\n", i)
		fmt.Printf("-------main %d\n", i)
		//go func2(ch, fmt.Sprintf("hello - %d\n", i))
	}

	//for i := 0; i < count; i++ {
	//	<-ch1
	//}
	for {

	}
}

func func2(ch chan int, data string) {

	fmt.Printf("%s", data)

	time.Sleep(time.Second * 3)

	client := ding_talk.NewClient("https://oapi.dingtalk.com/robot/send?access_token=15ecbce02525359580fddc1ac846873023224c3eb398d574977e7d5bf6dc5517")
	text := ding_talk.TextMessage{
		MsgType: ding_talk.Text,
		Text: ding_talk.TextData{
			Content: fmt.Sprintf("通知%s", data),
		},
	}
	client.Execute(text)
	ch <- 1
}

func func3() {
	for e := range ch2 {
		//data := <-ch

		fmt.Printf("%s", e)

		time.Sleep(time.Second * 2)

		//client := ding_talk.NewClient("https://oapi.dingtalk.com/robot/send?access_token=15ecbce02525359580fddc1ac846873023224c3eb398d574977e7d5bf6dc5517")
		//text := ding_talk.TextMessage{
		//	MsgType: ding_talk.Text,
		//	Text: ding_talk.TextData{
		//		Content: fmt.Sprintf("通知%s", e),
		//	},
		//}
		//client.Execute(text)
	}

}
