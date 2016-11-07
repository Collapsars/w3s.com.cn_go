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
	Type int
}

var (
	maxRoutineNum = 10
)

//移除教程多余的html 　含广告
func main() {
	fmt.Println("start")

	//ch := make(chan int , maxRoutineNum)

	//连接数据库
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=w3s_dev sslmode=disable password=root")
	utils.CheckErr(err)

	articles := []Article{}
	db.Where("type = 1").Find(&articles)

	count := 0

	for _,value := range articles {
		//fmt.Println(value.Type)

		doc , err := goquery.NewDocumentFromReader(strings.NewReader(value.Content))
		utils.CheckErr(err)

		doc.Find("h1").Each(func(i int , s *goquery.Selection) {
			//fmt.Println(s.Text())
			value.Title = s.Text()
		})




		//article := Article{};
		//article.Id =  value.Id



		//Content,err := doc.Html()

		//fmt.Println(value.Content)
		//fmt.Println(value.Id)
		///fmt.Println(time.Now().UnixNano())
		db.Save(&value)

		count ++


	}




	fmt.Println(count)


	fmt.Println("end")
	defer db.Close()


}

