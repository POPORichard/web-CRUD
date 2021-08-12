package model

import (
	"github.com/jinzhu/gorm"
	"reflect"
)

type Demo_order struct {

	gorm.Model
	Order_no string
	User_name string
	Amount float64
	Status string
	File_url string
}

type Demo_order_show struct {
	ID uint
	Order_no string
	User_name string
	Amount float64
	Status string
	File_url string
	CreateAt string
	UpdatedAt string

}
//判断结构体是否为空
func (dorder Demo_order)IsEmpty() bool{
	return reflect.DeepEqual(dorder, Demo_order{})
}

//将demo_order转为展示格式
func (order Demo_order)OrderToShow() Demo_order_show{
	show := Demo_order_show{}
	show.ID = order.Model.ID
	show.Order_no = order.Order_no
	show.User_name = order.User_name
	show.Amount = order.Amount
	show.Status = order.Status
	show.File_url = order.File_url
	show.CreateAt = order.Model.CreatedAt.Format("2006-01-02 15:04:05")
	show.UpdatedAt = order.Model.UpdatedAt.Format("2006-01-02 15:04:05")

	return show
}