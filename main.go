package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/malaohu/wechatmpbot/config"
	"github.com/malaohu/wechatmpbot/gtp"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func serveWechat(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.URL)
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID: config.LoadConfig().WxAppId,
		//AppSecret: "xxx",
		Token: config.LoadConfig().WxToken,
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	server := officialAccount.GetServer(req, rw)
	server.SetMessageHandler(handler)

	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	server.Send()
}

func handler(msg *message.MixMessage) *message.Reply {

	reply := "ChatGPT Bot 仅支持文字内容"
	if msg.MsgType == "text" {
		greply, _err := gtp.Completions(msg.Content)
		if _err != nil {
			log.Println("ERROR: ", _err)
		}
		reply = strings.TrimSpace(greply)
		reply = strings.Trim(reply, "\n")
	}
	text := message.NewText(reply)
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
}

func main() {
	port := ":" + config.LoadConfig().HttpPort
	http.HandleFunc("/", serveWechat)
	fmt.Println("wechat server listener at", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}
