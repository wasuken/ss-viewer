package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Url    string
	DatUrl string
	DbPath string
}

type SS struct {
	Id       string
	Title    string
	Contents string
}

var SSS []SS
var UrlConfig Config

func main() {
	_, err := toml.DecodeFile("./config.toml", &UrlConfig)
	if err != nil {
		panic(err)
	}
	SSS = listSS()
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	router.GET("/api/v1/ss/page", func(c *gin.Context) {
		c.JSON(200, gin.H{"page": ((len(SSS) - 1) / 100)})
	})
	router.GET("/api/v1/ss", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		page = page - 1

		if len(SSS) > page && 0 <= page {
			c.JSON(200, SSS[page*100:page*100+100])
		} else {
			c.JSON(200, SSS[(len(SSS)-101):(len(SSS)-1)])
		}
	})
	router.GET("/getDat", func(c *gin.Context) {
		res, err := http.Get(UrlConfig.Url)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		utfBody := transform.NewReader(bufio.NewReader(res.Body), japanese.ShiftJIS.NewDecoder())
		doc, err := goquery.NewDocumentFromReader(utfBody)
		if err != nil {
			panic(err)
		}
		var cnt int
		doc.Find("#threadlist td a").Each(func(index int, s *goquery.Selection) {
			if cnt > 1000 {
				return
			}
			href, _ := s.Attr("href")
			updateOrInsert(href, s.Text())
			cnt++
		})
	})
	router.Run(":8080")
}
func updateOrInsert(href, text string) {
	split := strings.Split(href, "/")
	id := split[len(split)-2]

	db, err := sql.Open("sqlite3", UrlConfig.DbPath)
	if err != nil {
		panic(err)
	}
	if containsKey(id) {
		fmt.Println(text)
		q := "insert into ss(title, id, contents) values(?,?,?)"
		stmt, _ := db.Prepare(q)
		_, _ = stmt.Exec(text, id, getDat(id))
		defer stmt.Close()
		time.Sleep(3 * time.Second)
		fmt.Println("inserted")
	} else {
		q := "update ss set contents = ? where id = ?"
		stmt, _ := db.Prepare(q)
		_, _ = stmt.Exec(getDat(id), id)
		defer stmt.Close()
		fmt.Println("updated")
	}
	time.Sleep(3 * time.Second)
	defer db.Close()
}
func containsKey(id string) bool {
	db, err := sql.Open("sqlite3", UrlConfig.DbPath)
	if err != nil {
		panic(err)
	}
	q := "select title from ss where id = ? Limit 1"
	defer db.Close()
	var title string
	if err := db.QueryRow(q, id).Scan(&title); err != nil {
		return true
	}
	if title != "" {
		return true
	} else {
		return false
	}
}
func getDat(id string) string {
	res, err := http.Get(UrlConfig.DatUrl + id + ".dat")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	utfBody := transform.NewReader(bufio.NewReader(res.Body), japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		panic(err)
	}
	html, _ := doc.Html()
	return html
}
func listSS() []SS {
	var sss []SS

	db, err := sql.Open("sqlite3", UrlConfig.DbPath)
	if err != nil {
		panic(err)
	}
	q := "select id,title,contents from ss"
	rows, _ := db.Query(q)
	defer rows.Close()
	for rows.Next() {
		var ss SS
		err = rows.Scan(&ss.Id, &ss.Title, &ss.Contents)
		if err != nil {
			panic(err)
		}
		sss = append(sss, ss)
	}
	defer db.Close()
	return sss
}
