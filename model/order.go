package model

import (
	"github.com/jinzhu/gorm"
	"reflect"
)
//DemoOrder 主要对象
type DemoOrder struct {

	gorm.Model
	OrderNo  string `gorm:"unique"`
	UserName string
	Amount   float64
	Status   string
	FileURL  string
}
//DemoOrderShow 用于展示的对象
type DemoOrderShow struct {
	ID        uint
	OrderNo   string
	UserName  string
	Amount    float64
	Status    string
	FileURL   string
	CreateAt  string
	UpdatedAt string

}
//IsEmpty 判断结构体是否为空
func (order DemoOrder)IsEmpty() bool{
	if &order == nil {
		return false
	}
	return reflect.DeepEqual(order, DemoOrder{})
}

//OrderToShow 将demo_order转为展示格式
func (order DemoOrder)OrderToShow() DemoOrderShow {
	show := DemoOrderShow{}
	show.ID = order.Model.ID
	show.OrderNo = order.OrderNo
	show.UserName = order.UserName
	show.Amount = order.Amount
	show.Status = order.Status
	show.FileURL = order.FileURL
	show.CreateAt = order.Model.CreatedAt.Format("2006-01-02 15:04:05")
	show.UpdatedAt = order.Model.UpdatedAt.Format("2006-01-02 15:04:05")

	return show
}