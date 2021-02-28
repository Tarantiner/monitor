package main

import "flag"

type dstInfo struct {
	Title  string
	Uid    string
	Link   string
	EPrice string
}

var t int
var checkt int
var items []*dstInfo

func init() {
	flag.IntVar(&t, "t", 30, "每隔t秒读取监控商品表")
	flag.IntVar(&checkt, "T", 300, "每隔t秒更新监控信息")
	flag.Parse()
}

func main() {
	// 从文本读取网页，解析出标题，价格线
	go func() {
		load()
	}()
	getCurrent()
}
