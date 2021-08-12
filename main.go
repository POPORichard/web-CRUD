package main

import (
	"fmt"
	"web_app/db"
	"web_app/router"
)

func main(){
	fmt.Println("start!")

	db.CreatePage()

	go db.AutoCheckEmptyURL()

	router.WebServer()



}


