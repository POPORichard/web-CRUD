package handler

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
	"web_app/db"
	"web_app/model"
)

var testdata = model.Demo_order{
	Model:     gorm.Model{},
	Order_no:  "test",
	User_name: "test",
	Amount:    0,
	Status:    "test",
	File_url:  "test",
}

func deleteTestData(){
	db := db.Link()
	db.Unscoped().Delete(&testdata)	//主键都没有怎么不会删错的？
	db.Close()
}

func TestNewData(t *testing.T) {
	NewData(&testdata)
	data := SearchByNo("test")
	if data != nil{
		if data.User_name != testdata.User_name ||
					data.Amount != testdata.Amount||
					data.Status != testdata.Status||
					data.File_url != testdata.File_url{
			fmt.Println(data)
			t.Errorf("data different")

		}
	}else {
		t.Errorf("data did't create, or data search failed!")
	}

	defer deleteTestData()
}


func TestAddFileURL(t *testing.T) {
	AddFileURL("test","test_URL")
	data := SearchByNo("test")
	if data != nil{
		if data.User_name != testdata.User_name ||
			data.Amount != testdata.Amount||
			data.Status != testdata.Status||
			data.File_url != "test_URL"{
			fmt.Println(data)
			t.Errorf("data different")

		}
	}else {
		t.Errorf("data did't create, or data search failed!")
	}

	defer deleteTestData()

}
