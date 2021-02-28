package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var err error
var ch = make(chan struct{}, 20)
var client http.Client
var lk sync.Mutex

type retInfo struct {
	dstInfo
	TM        string
	APrice    string
	Discount  string
	Status    string
	Baitiao   string
	NeedYuyue string
	YuyueTime string
}

type priceInfo struct {
	P        string `json:"p"`
	Discount string `json:"discount"`
}

type stockInfo struct {
	InfoStr string `json:"stockDesc"`
}

type baitiaoInfo struct {
	InfoStr string `json:"marketingText"`
}

type yuyueInfo struct {
	YuyueTime string `json:"yuyueTime"`
}

type goodsInfo struct {
	Price   priceInfo `json:"price"`
	Stock   stockInfo `json:"stockInfo"`
	Baitiao string    `json:"baitiaoPlanShowResVo"`
	Yuyue   yuyueInfo `json:"yuyueInfo"`
}

func download(item *dstInfo) {
	timeNow := time.Now().Format("2006/01/02 15:04:05")
	var req *http.Request
	u := "https://item-soa.jd.com/getWareBusiness"
	value := url.Values{
		"callback": []string{"jQuery2817335"},
		"skuId":    []string{item.Uid},
		//"cat": []string{"670,677,679"},
		"area": []string{"17_1381_50713_53756"},
		//"shopId": []string{"1000094211"},
		//"venderId": []string{"1000094211"},
		//"paramJson": []string{`{"platform2":"1","specialAttrStr":"p0pp1ppppppp2ppppppppppp","skuMarkStr":"00"}`},
		//"num": []string{"1"},
	}
	req, err = http.NewRequest("GET", u, nil)
	req.URL.RawQuery = value.Encode()
	if err != nil {
		log.Println("请求构造失败")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.190 Safari/537.36")
	rsp, err := client.Do(req)
	if err != nil {
		log.Println("请求失败")
		return
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println("read failed")
		return
	}

	reg, err := regexp.Compile(`jQuery\d+\((.*)\)`)
	if err != nil {
		log.Println("reg失败")
		return
	}
	ret := reg.FindSubmatch(b)

	var g goodsInfo
	err = json.Unmarshal(ret[1], &g)
	if err != nil {
		log.Println("json失败", err)
	}

	var bt baitiaoInfo
	if g.Baitiao != "" {
		json.Unmarshal([]byte(g.Baitiao), &bt)
	}

	status := "无货"
	if strings.Contains(g.Stock.InfoStr, "有货") {
		status = "有货"
	}

	needYuyue := "不需要"
	if g.Yuyue != (yuyueInfo{}) {
		needYuyue = "需要"
	}

	var info retInfo
	info.TM = timeNow
	info.Title = item.Title
	info.Link = item.Link
	info.EPrice = item.EPrice
	info.APrice = g.Price.P
	info.Discount = g.Price.Discount
	info.Status = status
	info.Baitiao = bt.InfoStr
	info.NeedYuyue = needYuyue
	info.YuyueTime = g.Yuyue.YuyueTime

	lk.Lock()
	fmt.Println(strings.Join([]string{info.TM, info.Title, info.APrice, info.EPrice, info.Discount, info.Status, info.Baitiao, info.NeedYuyue, info.YuyueTime, info.Link}, "  "))
	lk.Unlock()
}

//	s := fmt.Sprintf(`时间:%s
//单价:%s
//折扣:%s
//状态:%s
//白条:%s
//需要预约:%s
//预约时间:%s`, timeNow, info.EPrice, info., g.Price.Discount, status, bt.InfoStr, needYuyue, g.Yuyue.YuyueTime)
//	fmt.Println(s, len(retLis))
//	tt = append(tt, 2)
//}

func getCurrent() {
	for {
		fmt.Println(strings.Join([]string{"当前时间", "商品名", "价格", "期望价格", "折扣", "状态", "白条", "预约", "预约时间", "链接"}, "|"))
		for _, item := range items {
			go func(itm *dstInfo) {
				ch <- struct{}{}
				wg.Add(1)
				defer func() {
					<-ch
					wg.Done()
				}()
				download(itm)
			}(item)
		}
		wg.Wait()
		time.Sleep(time.Second * time.Duration(checkt))
	}

}
