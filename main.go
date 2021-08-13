package main

import (
	"fmt"
	"web_app/database"
	"web_app/router"
)

func main(){
	fmt.Println("start!")

	database.CreatePage()

	go database.AutoCheckEmptyURL()

	r := router.WebServer()
	r.Run(":8080")



}


