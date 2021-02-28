package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

func loadGoods() {
	var err error
	f, err := os.Open("goods.txt")
	if err != nil {
		log.Println("打开监控表失败", err)
		return
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	var s string
	for {
		s, err = rd.ReadString('\n')
		data := strings.SplitN(s, "|", 4)
		if len(data) == 4 {
			items = append(items, &dstInfo{
				Title:  data[0],
				Uid:    data[1],
				EPrice: data[2],
				Link:   data[3],
			})
		}
	}
}

func load() {
	// 加载文本
	for {
		items = []*dstInfo{}
		loadGoods()
		time.Sleep(time.Second * time.Duration(t))
	}
}
