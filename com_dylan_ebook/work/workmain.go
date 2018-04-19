package work

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strconv"
	"strings"
	"runtime"
)

//我只要数据 至于它是什么类型
//分析

const ebookUrl = "http://www.woaidu.org/book_1.html"

type EbookData struct {
	DataUrl  string
	ImageUrl string
	Title    string
	Content  string
	Time     string
	Author   string
	Download string
}

type DownLoadInfo struct {
	DownLoadUrl string
	Downfrom    string
	Updatetime  string
	Progress    string
}

var operater *Operater

func work() {

}

type person struct {
	AGE    int    `bson:"age"`
	NAME   string `bson:"name"`
	HEIGHT int    `bson:"height"`
}

func Start() {
	runtime.GOMAXPROCS(2)

	jobs := make(chan string, 100)
	results := make(chan int, 100000)

	operater = new(Operater)
	operater.Dbname = "ebook";
	operater.Document = "ebookdata";
	err := operater.Connect()
	if err != nil {
		fmt.Println(err)
	}
	urls, count := getUrlData(19900, 100000)

	defer operater.Close()

	// 3是缓存池大小
	for w := 1; w <= 3; w++ {
		go worker(jobs, results)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)
	// Finally we collect all the results of the work.
	for a := 1; a <= count; a++ {
		<-results
	}
}

func worker(jobs chan string, results chan<- int) {

	for url := range jobs {
		ebookData, err := getData(url)
		if err != nil {
			fmt.Println(err)
		} else {
			ebookData.DataUrl = url
			//ebookData:=new(EbookData)
			fmt.Println(url)
			err = operater.Insert(ebookData)
			fmt.Println(err)
		}
		results <- 1
	}

}

func getUrlData(start int, end int) ([]string, int) {
	urls := make([]string, 0)
	for i := start; i < end; i++ {
		s := "http://www.woaidu.org/book_" + strconv.Itoa(i) + ".html"
		urls = append(urls, s)
	}
	return urls, end - start
}

func parseImageUrl(reader io.Reader) (*EbookData, error) {
	d := new(EbookData)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Println(err)
		return d, err
	}

	d.ImageUrl, _ = doc.Find(".hong img").Attr("src")
	d.Title = doc.Find(".zizida").Text()
	d.Author = doc.Find(".xiaoxiao").Text()
	d.Content = doc.Find(".lili").Text()
	d.Time = doc.Find(".jiewei").Text()
	downloadinfos := make([]DownLoadInfo, 0)

	doc.Find("div.xiazai_xiao").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			//fmt.Println(s.Html())
			url, exist := s.Find(".pcdownload").Find("a").Attr("href")
			if exist {
				if strings.HasSuffix(url, "zip") || strings.HasSuffix(url, "rar") {
					downinfo := new(DownLoadInfo)
					downinfo.DownLoadUrl = url
					s.Find(".ziziziz").Each(func(i int, s1 *goquery.Selection) {
						if i == 1 {
							downinfo.Progress = s1.Text()
						} else if i == 2 {
							downinfo.Updatetime = s1.Text()
						} else {
							downinfo.Downfrom = s1.Text()
						}
					})
					downloadinfos = append(downloadinfos, *downinfo)
				}
			}
		}
	})

	info, err := json.Marshal(downloadinfos)
	d.Download = string(info)

	return d, nil
}

func getData(url string) (*EbookData, error) {
	req := buildRequest(url)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("获取 io ")
	}
	return parseImageUrl(io.Reader(resp.Body))

}

func buildRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")
	//req.Header.Set("Cookie", )
	return req
}
