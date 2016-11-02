package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"w3s.com.cn/utils"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
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
	db.Where("type = 2").Find(&articles)

	count := 0

	for _,value := range articles {
		//fmt.Println(value.Type)



		doc , err := goquery.NewDocumentFromReader(strings.NewReader(value.Content))
		utils.CheckErr(err)







		//移除上面广告
		doc.Find(".ad").Each(func(i int , s *goquery.Selection) {
			//fmt.Println(s.Html())
			s.Remove()
		})

		//移除上面广告
		doc.Find("label.pull-right").Each(func(i int , s *goquery.Selection) {
			//fmt.Println(s.Html())
			s.Remove()
		})


		//移除div display:none;
		doc.Find("div").Each(func(i int , s *goquery.Selection) {
			stylet , _ :=s.Attr("style")
			//fmt.Println(stylet)
			if(stylet == "display:none;"){
				s.Remove()
			}

			//fmt.Println(s.Html())
			//
			//s.Remove()
		})





		r, _ := regexp.Compile("runcode")
		//fmt.Println(scriptText)
		//http://tool.runoob.com/compile.php

		scriptText, _:= doc.Html()






		if (r.MatchString(scriptText)){
			//fmt.Println(scriptText)
			rr, _ := regexp.Compile("http://tool.runoob.com/compile.php")
			value.Content = string(rr.ReplaceAll([]byte(scriptText),[]byte("/html/run")))
		}else {
			value.Content = scriptText
		}



		////去掉
		//doc.Find("a").Each(func(i int , s *goquery.Selection) {
		//	link,_ := s.Attr("href")
		//	r, _ := regexp.Compile("^http://www.runoob.com/")
		//	if r.MatchString(link){
		//		//fmt.Println(link)
		//		path := string(r.ReplaceAll([]byte(link),[]byte("/")))
		//		s.SetAttr("href",path)
		//		//fmt.Println(path)
		//	}
		//})

		//article := Article{};
		//article.Id =  value.Id

		//value.Content = scriptText

		//Content,err := doc.Html()

		//fmt.Println(value.Content)
		//fmt.Println(value.Id)
		fmt.Println(time.Now().UnixNano())
		//fmt.Println(value.Content)
		//fmt.Println(value.Content)
		db.Save(&value)

		count ++
		/*
		ch <- 1
       		go func() {
			doc , err := goquery.NewDocumentFromReader(strings.NewReader(value.Content))
			utils.CheckErr(err)

			//移除右边
			doc.Find(".right-column").Each(func(i int , s *goquery.Selection) {
				//fmt.Println(s.Html())
				s.Remove()
			})

			/*
			//移除广告
			doc.Find("ins").Each(func(i int , s *goquery.Selection) {
				s.Remove()
			})
			//移除google script广告
			doc.Find("script").Each(func(i int , s *goquery.Selection) {
				s.Remove()
			})

			//移除google script广告
			doc.Find("script").Each(func(i int , s *goquery.Selection) {
				s.Remove()
			})

			//移除试一试菜鸟工具
			doc.Find("#RightPane .pull-right a").Each(func(i int , s *goquery.Selection) {
				s.Remove()
			})




			//article := Article{};
			//article.Id =  value.Id

			value.Content,err = doc.Html()

			//Content,err := doc.Html()

			//fmt.Println(value.Content)
			fmt.Println(value.Id)
			fmt.Println(time.Now().UnixNano())
			db.Save(&value)
			//count ++
			<- ch
		}()

	*/
		//fmt.Println(doc.Html())

	}




	fmt.Println(count)


	fmt.Println("end")
	defer db.Close()


}

