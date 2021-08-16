package model

import (
	"fmt"
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
	ID        uint		`json:"id"`
	OrderNo   string	`json:"order_no"`
	UserName  string	`json:"user_name"`
	Amount    float64	`json:"amount"`
	Status    string	`json:"data_status"`
	FileURL   string	`json:"file_url"`
	CreateAt  string	`json:"create_at"`
	UpdatedAt string	`json:"updated_at"`

}
//IsEmpty 判断结构体是否为空
func (order *DemoOrder)IsEmpty() bool{
	if order == nil {
		return true
	}
	fmt.Println("The empty data is:",order)
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