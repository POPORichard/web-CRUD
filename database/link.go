package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"	//添加gorm对mysql语句的支持
	"web_app/model"
)

//Link 连接数据库
func Link()(db *gorm.DB){
	db,err := gorm.Open("mysql","web_app:123456@(localhost)/web_app?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil{
		fmt.Println("Link to database error :",err)
	}

	return
}

//CreatePage 根据结构体自动创建表
func CreatePage(){
	var model model.DemoOrder

	db := Link()
	db.AutoMigrate(model)
	db.Close()
}