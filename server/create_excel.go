package server

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"strconv"
	"sync"
	"web_app/handler"
	"web_app/model"
)

func WriteToExcel(data []model.DemoOrder, waitGroup *sync.WaitGroup)error{
	defer waitGroup.Done()

	fmt.Println("start create exlse!")
	if data ==nil{
		data,_ = handler.GetAllData()
	}

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	row.SetHeightCM(2)
	cell = row.AddCell()
	cell.Value = "id"
	cell = row.AddCell()
	cell.Value = "create_time"
	cell = row.AddCell()
	cell.Value = "update_time"
	cell = row.AddCell()
	cell.Value = "order_no"
	cell = row.AddCell()
	cell.Value = "user_name"
	cell = row.AddCell()
	cell.Value = "amount"
	cell = row.AddCell()
	cell.Value = "status"
	cell = row.AddCell()
	cell.Value = "file_url"

	for i := range data{
		row = sheet.AddRow()
		row.SetHeightCM(1)
		cell = row.AddCell()
		cell.Value = strconv.Itoa(int(data[i].ID))
		cell = row.AddCell()
		cell.Value = data[i].CreatedAt.Format("2006-01-02 15:04:05")
		cell = row.AddCell()
		cell.Value = data[i].UpdatedAt.Format("2006-01-02 15:04:05")
		cell = row.AddCell()
		cell.Value = data[i].OrderNo
		cell = row.AddCell()
		cell.Value = data[i].UserName
		cell = row.AddCell()
		cell.Value = strconv.FormatFloat(data[i].Amount,'e',-1,64)
		cell = row.AddCell()
		cell.Value = data[i].Status
		cell = row.AddCell()
		cell.Value = data[i].FileURL

	}


	err = file.Save("../list.xlsx")
	if err != nil {
		fmt.Println("Create xlsx error:",err.Error())
	}else{
		fmt.Println("Finish create xlsx!")
	}


	return nil
}
