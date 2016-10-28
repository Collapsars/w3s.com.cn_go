package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"w3s.com.cn/utils"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/PuerkitoBio/goquery"
	"strings"
)


//model
type Article struct {
	Id int
	Title string
	Description string
	Path string
	Content string
}

func main() {
	fmt.Println("start")

	//连接数据库
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=w3s_dev sslmode=disable password=root")
	utils.CheckErr(err)

	articles := []Article{}
	db.Find(&articles)

	//count := 0

	for _,value := range articles {

		doc , err := goquery.NewDocumentFromReader(strings.NewReader(value.Content))
		utils.CheckErr(err)
		//找到试一试按钮
		doc.Find(".middle-column .tryitbtn").Each(func(i int , s *goquery.Selection) {
			tryLink , _ := s.Attr("href")
			spiderTryLink(tryLink,db)

			//fmt.Println(tryLink)
			//html , err := goquery.NewDocument("http://www.runoob.com"+tryLink)
			//fmt.Println(s.Attr("href"))
			//count ++
		})

		doc.Find(".middle-column .showbtn").Each(func(i int , s *goquery.Selection) {
			tryLink , _ := s.Attr("href")
			spiderTryLink(tryLink,db)
			//fmt.Println(tryLink)
			//html , err := goquery.NewDocument("http://www.runoob.com"+tryLink)
			//fmt.Println(s.Attr("href"))
			//count ++
		})




	}



	//fmt.Println(count)

	defer db.Close()


}

//下载连接
func spiderTryLink(url string,db *gorm.DB)  {



	//先判断
	articlee :=Article{}
	db.Where("path = ?", url).First(&articlee)
	if  articlee.Id == 0{

		html , err := goquery.NewDocument("http://www.runoob.com"+url)
		utils.CheckErr(err)
		//找到.panel-body
		html.Find(".panel-body").Each(func(i int , s *goquery.Selection) {
			tryHtml , _ := s.Html()
			//fmt.Println(tryHtml)




			artile := Article{
				Path:url,
				Content:tryHtml,
			}

			db.Create(&artile)




		})



	}






}


