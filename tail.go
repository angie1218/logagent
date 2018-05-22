package main

import (
	"time"

	"github.com/apex/log"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
)

type CollectConf struct {
	LogPath  string
	Topic    string
	ChanSize int
}

type TailObj struct {
	tail *tail.Tail
	conf CollectConf
}

type TailObjMgr struct {
	tailsObjs []*TailObj
	msgChan   chan *TextMsg
}

type TextMsg struct {
	Msg   string
	Topic string
}

var (
	tailObjMgr *TailObjMgr
)

func GetOneLine() (msg *TextMsg) {

	msg = <-tailObjMgr.msgChan
	return
}

func InitTail(confs []CollectConf, chanSize int) (err error) {

	// 容错处理
	if len(confs) == 0 {
		logs.Error("invaild config for log collect, conf:%v", confs)
		return
	}
	// 初始化管道
	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *TextMsg, chanSize),
	}
	for _, v := range confs {

		obj := &TailObj{
			conf: v,
		}

		tails, errTail := tail.TailFile(v.LogPath, tail.Config{
			ReOpen:    true,
			Follow:    true,
			MustExist: false,
			Poll:      true,
		})

		if errTail != nil {

			log.Errorf("tailf occurs errors, error: %v", err)

			return
		}

		obj.tail = tails
		tailObjMgr.tailsObjs = append(tailObjMgr.tailsObjs, obj)

		go ReadFromTail(obj)

	}

	return
}

func ReadFromTail(tailObj *TailObj) {

	for true {
		line, ok := <-tailObj.tail.Lines
		if !ok {
			logs.Warn("tail file close reopen, filename:%s\n", tailObj.tail.Filename)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		textMsg := &TextMsg{
			Msg:   line.Text,
			Topic: tailObj.conf.Topic,
		}

		//tailObjMgr := &TailObjMgr{}
		tailObjMgr.msgChan <- textMsg
		return
	}
}