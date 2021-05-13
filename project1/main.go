package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type FormData struct {
	IpAddress  string `json:"ipaddress"`
	IpLocation string `json:"location"`
	Image      string `json:"image"`
}

func main() {
	var router = gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))
	router.POST("/save_form", FormHandler)
/* 	router.GET("/get_image", ImageHandler) */
	router.Run(":8080")
}
func FormHandler(c *gin.Context) {
	var data FormData
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
	}

	db := DBconnect()
	defer db.Close()

	_, err = db.Exec("insert into public.data(ipaddress,location,image) values($1,$2,$3)", data.IpAddress, data.IpLocation, data.Image)
	if err != nil {
		fmt.Println("here")
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "data saved",
	})
}
/* func ImageHandler(c *gin.Context) {
	var data FormData
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
	}

	db := DBconnect()
	defer db.Close()
	var image string
	res, err := db.Query("select image from public.data where ipaddress=$1", data.IpAddress)
	if err != nil {
		fmt.Println(err)
	}
	for res.Next() {
		err = res.Scan(&image)
		if err != nil {
			fmt.Println(err)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"image":  image,
	})
} */
func DBconnect() *sql.DB {

	db, err := sql.Open("postgres", "user=postgres password=root dbname=test sslmode=disable")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(db)
	// close database

	// check db
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")
	return db
}
