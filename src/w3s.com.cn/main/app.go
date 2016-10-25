package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"w3s.com.cn/utils"
	"regexp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Article struct {
	Id int
	Title string
	Description string
	Path string
	Content string
}

var c chan int

func main() {
	for i := 0 ; i < 4 ; i++  {
		go Scrape("http://www.runoob.com/")
	}



	fmt.Println(<- c)


}

func Scrape(url string , ) () {


	doc , err := goquery.NewDocument(url)
	utils.CheckErr(err)

	home := make([]string,0,200)
	doc.Find(".item-top").Each(func(i int , s *goquery.Selection) {
		//title := s.Find("h2").Text()
		//fmt.Println(s.Attr("href"))
		title_link, _ := s.Attr("href")
		home = append(home,title_link)
	})
	//主页
	//home := make([]string,0,200) // 首页导航标签切片
	//home = append(home,"http://www.runoob.com/cprogramming/c-tutorial.html")

	page := make([]string,0,200) //每个导航页面切片
	for _ , value := range home {

		//fmt.Println(value)
		html , err := goquery.NewDocument(value)
		utils.CheckErr(err)

		html.Find("#leftcolumn a").Each(func(i int , s *goquery.Selection) {
			//fmt.Println(s.Attr("href"))
			a_name , _ := s.Attr("href")
			r, _ := regexp.Compile("^http://www.runoob.com")
			r1, _ := regexp.Compile("^/")
			if !r.MatchString(a_name) {
				if r1.MatchString(a_name) {
					a_name = "http://www.runoob.com"+a_name
				}else {
					a_name = "http://www.runoob.com/"+a_name
				}

				//fmt.Println(a_name)
			}
			page = append(page , a_name)


		})
	}

	//fmt.Println(page)

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=w3s_dev sslmode=disable password=root")
	utils.CheckErr(err)



	//pagee := make([]string,0,200) //每个导航页面切片
	for _ , value := range page {

		 loop(value,db)

	}


	//defer db.Close()


	c <- 1
	//fmt.Println("end")


}

//当前链接里面的所有链接
func loop(str string,db *gorm.DB)  {
	//fmt.Println(value)
	html , err := goquery.NewDocument(str)
	utils.CheckErr(err)

	html.Find("#leftcolumn a").Each(func(i int , s *goquery.Selection) {
		//fmt.Println(s.Attr("href"))
		a_name , _ := s.Attr("href")
		r, _ := regexp.Compile("^http://www.runoob.com")
		r1, _ := regexp.Compile("^/")
		if !r.MatchString(a_name) {
			if r1.MatchString(a_name) {
				a_name = "http://www.runoob.com"+a_name
			} else {
				a_name = "http://www.runoob.com/"+a_name
			}

			//fmt.Println(a_name)
		}



		//去掉http://www.runoob.com/
		re, _ := regexp.Compile("^http://www.runoob.com/")
		a_name = string(re.ReplaceAll([]byte(a_name),[]byte("/")))

		articlee :=Article{}
		db.Where("path = ?", a_name).First(&articlee)

		//fmt.Println(articlee.Id)

		if  articlee.Id == 0{
			article , err := goquery.NewDocument("http://www.runoob.com"+a_name)
			utils.CheckErr(err)

			content , err := article.Find(".main").Html()
			utils.CheckErr(err)


			artile := Article{
				Path:a_name,
				Content:content,
			}

			db.Create(&artile)
		}

		//loop("http://www.runoob.com"+a_name,db)
	})





}

