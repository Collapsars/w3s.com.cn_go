package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"w3s.com.cn/utils"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
)


//model
type Article struct {
	Id int
	Title string
	Description string
	Path string
	Content string
	Type int
}



var (
	maxRoutineNum = 10
)


//爬取每个教程中试一试到数据库
func main() {
	fmt.Println("start")
	ch := make(chan int , maxRoutineNum)

	//连接数据库
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=w3s_dev sslmode=disable password=root")
	utils.CheckErr(err)

	articles := []Article{}
	db.Find(&articles)

	//count := 0

	//fmt.Println("2")
	for _,value := range articles {
		//fmt.Println("1")

		doc , err := goquery.NewDocumentFromReader(strings.NewReader(value.Content))
		utils.CheckErr(err)
		//找到试一试按钮
		doc.Find(".middle-column .tryitbtn").Each(func(i int , s *goquery.Selection) {
			tryLink , _ := s.Attr("href")
			//fmt.Println(tryLink)
			ch <- 1
			go spiderTryLink(tryLink,db,ch)

			//fmt.Println(tryLink)
			//html , err := goquery.NewDocument("http://www.runoob.com"+tryLink)
			//fmt.Println(s.Attr("href"))
			//count ++
		})

		doc.Find(".middle-column .showbtn").Each(func(i int , s *goquery.Selection) {
			tryLink , _ := s.Attr("href")
			ch <- 1
			go spiderTryLink(tryLink,db,ch)
			//fmt.Println(tryLink)
			//html , err := goquery.NewDocument("http://www.runoob.com"+tryLink)
			//fmt.Println(s.Attr("href"))
			//count ++
		})




	}



	defer db.Close()
	fmt.Println("end")


}

//下载连接
func spiderTryLink(url string,db *gorm.DB,ch chan int)  {



	//先判断
	articlee :=Article{}
	db.Where("path = ?", url).First(&articlee)

	if  articlee.Id == 0{
		//fmt.Println(articlee.Id)
		html , err := goquery.NewDocument("http://www.runoob.com"+url)
		utils.CheckErr(err)
		fmt.Println("http://www.runoob.com"+url)
		//找到.panel-body
		html.Find("body").Each(func(i int , s *goquery.Selection) {
			tryHtml , _ := s.Html()
			//fmt.Println(tryHtml)

			artile := Article{
				Path:url,
				Content:tryHtml,
				Type:2,
			}

			db.Create(&artile)


		})
	}
	fmt.Println(time.Now().UnixNano())
	<- ch




}


