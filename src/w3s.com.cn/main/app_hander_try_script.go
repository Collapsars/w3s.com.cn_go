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
		//html 的script

		doc.Find(".middle-column .tryitbtn").Each(func(i int , s *goquery.Selection) {
			tryLink , _ := s.Attr("href")
			spiderTryLinkk(tryLink,db)
		})



		doc.Find(".middle-column .showbtn").Each(func(i int , s *goquery.Selection) {
			tryLink , _ := s.Attr("href")
			spiderTryLinkk(tryLink,db)
		})





	}



	//fmt.Println(count)

	defer db.Close()


}

//下载连接
func spiderTryLinkk(url string,db *gorm.DB)  {



	//先判断
	articlee :=Article{}
	db.Where("path = ?", url).First(&articlee)
	if  articlee.Id > 0{

		html , err := goquery.NewDocument("http://www.runoob.com"+url)
		utils.CheckErr(err)
		//找到.panel-body
		html.Find("body script").Each(func(i int , s *goquery.Selection) {

			//fmt.Println(e.Text())
			scriptText:= s.Text()
			//r, _ := regexp.Compile("^var\\s*seditor")
			//r, _ := regexp.Compile("//\\s*Define")
			r, _ := regexp.Compile("runcode")
			//fmt.Println(scriptText)
			//http://tool.runoob.com/compile.php
			if (r.MatchString(scriptText)){
				//fmt.Println(scriptText)
				rr, _ := regexp.Compile("http://tool.runoob.com/compile.php")
				scriptText = string(rr.ReplaceAll([]byte(scriptText),[]byte("/html/run")))
				fmt.Println(scriptText)
				articlee.Description = scriptText
				db.Save(&articlee)

			}


		})



	}






}


