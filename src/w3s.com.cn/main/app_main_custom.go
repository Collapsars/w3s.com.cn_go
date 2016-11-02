package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"w3s.com.cn/utils"
	"time"
	"regexp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strings"
)

type Article struct {
	Id int
	Title string
	Description string
	Path string
	Content string
}

var (
	maxRoutineNum = 10
)

func main()  {
	fmt.Println("start")
	ch := make(chan int , maxRoutineNum)
	for  i := 0; i< 5 ; i++ {
		spiderMainn(ch)
	}


	fmt.Println("end")
}

func spiderMainn(ch chan int)  {


	//连接数据库
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=w3s_dev sslmode=disable password=root")
	utils.CheckErr(err)

	//count := 0
	home := []string{
		"/memcached/memcached-tutorial.html",
		"/w3cnote/android-tutorial-intro.html",

	}


	///fmt.Println(count)


	for _  , value :=range home{
		html , err := goquery.NewDocument("http://www.runoob.com"+value)
		utils.CheckErr(err)

		//fmt.Println(html)

		pre := strings.Split(value,"/")
		//fmt.Println(pre[1])

		html.Find("#leftcolumn a").Each(func(i int , s *goquery.Selection) {
			//fmt.Println(s.Attr("href"))
			a_name , _ := s.Attr("href")
			r, _ := regexp.Compile("^http://www.runoob.com")
			r1, _ := regexp.Compile("^/")
			r2, _ := regexp.Compile("^http://www.w3cschool.cc")
			if !r.MatchString(a_name) {
				if r1.MatchString(a_name) {
					a_name = "http://www.runoob.com"+a_name
				}else if r2.MatchString(a_name){

				}else {
					//fmt.Println(a_name)
					//writeFile(a_name+"\n")
					a_name = "http://www.runoob.com/"+pre[1]+"/"+a_name
				}
			}
			//fmt.Println(a_name)
			ch <- 1
			go handerUrll(a_name,db,ch)

		})
	}






}



//处理url到数据库
func handerUrll(str string,db *gorm.DB,ch chan int)  {

	re, _ := regexp.Compile("^http://www.runoob.com/")
	ree, _ := regexp.Compile("^http://www.w3cschool.cc/")

	var url string
	if  re.MatchString(str){
		url = string(re.ReplaceAll([]byte(str),[]byte("/")))
	}else {
		url = string(ree.ReplaceAll([]byte(str),[]byte("/")))
	}

	articlee :=Article{}
	db.Where("path = ?", url).First(&articlee)

	if articlee.Id == 0 {
		article , err := goquery.NewDocument(str)
		utils.CheckErr(err)

		content , err := article.Find(".main").Html()
		utils.CheckErr(err)


		artile := Article{
			Path:url,
			Content:content,
		}

		db.Create(&artile)
	}
	fmt.Println(time.Now().UnixNano())
	<- ch
}
