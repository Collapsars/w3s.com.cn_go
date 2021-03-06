package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"w3s.com.cn/utils"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"regexp"
	"time"
)


//model
type Article struct {
	Id int
	Title string
	Description string
	Path string
	Content string
}

var (
	maxRoutineNum = 10 ;
)

//去年http://www.runoob.com/
func main() {
	fmt.Println("start")

	ch := make(chan int ,maxRoutineNum)

	//连接数据库
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=w3s_dev sslmode=disable password=root")
	utils.CheckErr(err)

	articles := []Article{}
	db.Find(&articles)

	for _,value := range articles {

		ch <- 1
		go func() {
			doc , err := goquery.NewDocumentFromReader(strings.NewReader(value.Content))
			utils.CheckErr(err)

			//去掉
			doc.Find("a").Each(func(i int , s *goquery.Selection) {
				link,_ := s.Attr("href")
				r, _ := regexp.Compile("^http://www.runoob.com/")
				if r.MatchString(link){
					//fmt.Println(link)
					path := string(r.ReplaceAll([]byte(link),[]byte("/")))
					s.SetAttr("href",path)
					//fmt.Println(path)
				}
			})


			//article := Article{};
			//article.Id =  value.Id

			value.Content,err = doc.Html()

			//Content,err := doc.Html()

			//fmt.Println(article)
			db.Save(&value)
			fmt.Println(time.Now().UnixNano())
			<- ch
		}()



		fmt.Println("end")
		//fmt.Println(doc.Html())

	}










	defer db.Close()


}

