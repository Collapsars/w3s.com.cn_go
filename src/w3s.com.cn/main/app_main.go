package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"w3s.com.cn/utils"
	"time"
	"regexp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
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
	for  i := 0; i< 50 ; i++ {
		spiderMain("http://www.runoob.com/sitemap",ch)
	}
	//spiderMain("http://www.runoob.com/sitemap",ch)

	fmt.Println("end")
}

func spiderMain(url string,ch chan int)  {
	//db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=w3s_dev sslmode=disable password=root")
	//utils.CheckErr(err)

	//连接数据库
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=w3s_dev sslmode=disable password=root")
	utils.CheckErr(err)


	doc , err := goquery.NewDocument(url)
	utils.CheckErr(err)
	//count := 0
	home := make([]string,0,200)
	doc.Find(".middle-column .sitemap ul li a").Each(func(i int , s *goquery.Selection) {
		main_link, _ := s.Attr("href")
		r, _ := regexp.Compile("^http://www.runoob.com/")
		main_link = string(r.ReplaceAll([]byte(main_link),[]byte("/")))
		//fmt.Println(main_link)
		home = append(home,main_link)
		//count ++
	})

	///fmt.Println(count)


	for _  , value :=range home{
		html , err := goquery.NewDocument("http://www.runoob.com"+value)
		utils.CheckErr(err)
		//fmt.Println(value)
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
			go handerUrl(a_name,db,ch)

		})
	}






}

//单个写入  文件手动处理
func writeFile(str string)  {
	f,err:=os.OpenFile("test.txt",os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		panic(err)
	}
	f.WriteString(str)
	defer f.Close()
}

//处理url到数据库
func handerUrl(str string,db *gorm.DB,ch chan int)  {

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
