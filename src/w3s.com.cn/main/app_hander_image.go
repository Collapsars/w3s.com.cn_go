package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"w3s.com.cn/utils"
	"github.com/PuerkitoBio/goquery"
	"strings"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"regexp"
	"net/http"
	"os"
	"io"
	"time"
	"path/filepath"
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
	maxRoutineNum = 10
)

//下载图片
func main()  {
	fmt.Println("start")

	 t1 := time.Now().UnixNano()
	//fmt.Println(time.Now()) 197530895100

	ch := make(chan int, maxRoutineNum)



	//连接数据库
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=w3s_dev sslmode=disable password=root")
	utils.CheckErr(err)

	articles := []Article{}
	db.Find(&articles)

	for _,value := range articles {

		doc , err := goquery.NewDocumentFromReader(strings.NewReader(value.Content))
		utils.CheckErr(err)
		//分割图片路径
		doc.Find("img").Each(func(i int , s *goquery.Selection) {
			//fmt.Println(s.Attr("src"))
			path , _ := s.Attr("src")

			re, _ := regexp.Compile("^http://www.runoob.com/")
			path = string(re.ReplaceAll([]byte(path),[]byte("/")))
			//fmt.Println(path)

			ch <- 1

			go downloadImg(path,ch)


		})

	}



	defer db.Close()

	t2 := time.Now().UnixNano()

	fmt.Println(t2-t1)

	fmt.Println("end")
}

func downloadImg(url string,ch chan int)  {


	r, _ := regexp.Compile("^/")
	//fmt.Println(url)
	if r.MatchString(url){
		res, err := http.Get("http://www.runoob.com"+url)
		utils.CheckErr(err)
		defer res.Body.Close()

		//r, _ := regexp.Compile("/")
		//url = string(r.ReplaceAll([]byte(url),[]byte("\\")))

		url = filepath.FromSlash(url)

		//fmt.Println(url)





		fileEixt ,_ := PathExists(filepath.Dir("E:\\project\\w3s.com.cn\\app\\assets\\test"+url))

		//fmt.Println(fileEixt)



		if !fileEixt {

			//fmt.Println("path:"+"E:\\project\\w3s.com.cn\\app\\assets\\test"+url)
			foldPath := filepath.Dir("E:\\project\\w3s.com.cn\\app\\assets\\test"+url)
			//fmt.Println("dir:"+foldPath)

			err := os.MkdirAll(foldPath,os.ModePerm)
			//fmt.Println("dir"+path.Dir("E:\\project\\w3s.com.cn\\app\\assets\\test.png"+url))

			utils.CheckErr(err)
			//os.Remove("E:\\project\\w3s.com.cn\\app\\assets\\test"+url)

		}


		file , err := os.Create("E:\\project\\w3s.com.cn\\app\\assets\\test"+url)
		utils.CheckErr(err)
		io.Copy(file,res.Body)

		fmt.Println(time.Now().UnixNano())


	}

	<-ch




}


func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}