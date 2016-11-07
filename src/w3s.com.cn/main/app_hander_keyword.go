package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"w3s.com.cn/utils"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"regexp"
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
	db.Find(&articles)

	count := 0

	for _,value := range articles {
		//fmt.Println(value.Type)

		doc , err := goquery.NewDocumentFromReader(strings.NewReader(value.Content))
		utils.CheckErr(err)

		r, _ := regexp.Compile("菜鸟教程")
		rr, _ := regexp.Compile("www.runoob.com")
		rrr, _ := regexp.Compile("(runoob.com)")


		src , _ := doc.Html()

		//菜鸟教程 - >www.w3s.com.cn
		//www.runoob.com - >www.w3s.com.cn
		src = r.ReplaceAllString(src,"www.w3s.com.cn")
		src = rr.ReplaceAllString(src,"www.w3s.com.cn")
		src = rrr.ReplaceAllString(src,"")
		//fmt.Println(src)
		//if r.MatchString(doc.Text()) {
		//	fmt.Println("true")
		//}



		//article := Article{};
		//article.Id =  value.Id

		value.Content = src

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

