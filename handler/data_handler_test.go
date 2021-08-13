package handler

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
	"web_app/database"
	"web_app/model"
)

var testdata = model.DemoOrder{
	Model:    gorm.Model{},
	OrderNo:  "test",
	UserName: "test",
	Amount:   0,
	Status:   "test",
	FileUrl:  "test",
}

func deleteTestData(){
	db := database.Link()
	db.Unscoped().Delete(&testdata)
	db.Close()
}

func TestNewData(t *testing.T) {
	NewData(&testdata)
	data := SearchByNo("test")
	if data != nil{
		if data.UserName != testdata.UserName ||
					data.Amount != testdata.Amount||
					data.Status != testdata.Status||
					data.FileUrl != testdata.FileUrl {
			t.Errorf("data different")

		}
	}else {
		t.Errorf("data did't create, or data search failed!")
	}

	defer deleteTestData()
}


func TestAddFileURL(t *testing.T) {
	NewData(&testdata)

	AddFileURL("test","test_URL")
	data := SearchByNo("test")
	if data != nil{
		if data.UserName != testdata.UserName ||
			data.Amount != testdata.Amount||
			data.Status != testdata.Status||
			data.FileUrl != "test_URL"{
			t.Errorf("data different")

		}
	}else {
		t.Errorf("data did't create, or data search failed!")
	}

	defer deleteTestData()
}


func TestSearch(t *testing.T) {

	NewData(&testdata)
	defer deleteTestData()

	pageSize := GetTotalNumber()/3
	data,err := CutData(3,pageSize,nil)
	if err != nil{
		t.Fatal("CutData fatal :",err)
	}

	searchData,err := Search("test",data)

	if err != nil{
		t.Errorf("Search error :%v",err )
	}

	for i := range searchData{
		fmt.Println("searchdata no:",i,":",searchData[i].UserName)
	}
	fmt.Println("please check!")
}

func TestSequenceByAmount(t *testing.T) {
	pageSize := GetTotalNumber()/3
	data,err := CutData(2,pageSize,nil)
	if err != nil{
		t.Fatal("CutData fatal :",err)
	}

	sequenceData,err := Sequence("amount", data)
	if err != nil{
		t.Errorf("Sequence error:%v",err)
		return
	}

	var tmp = sequenceData[0].Amount
	fmt.Println("Sort by amount is:")
	 for i := range sequenceData{
	 	amount := sequenceData[i].Amount
	 	fmt.Printf("%v>",amount)
	 	if tmp < amount{
	 		t.Errorf("sequenceData error around %v",sequenceData[i].Amount)
		}
		tmp = amount

	 }
}

func TestSequenceByTime(t *testing.T) {
	pageSize := GetTotalNumber()/3
	data,err := CutData(2,pageSize,nil)
	if err != nil{
		t.Fatal("CutData fatal :",err)
	}

	sequenceData,err := Sequence("time", data)
	if err != nil{
		t.Errorf("Sequence error:%v",err)
		return
	}

	var tmp = sequenceData[0].UpdatedAt
	fmt.Println("Sort by time is:")
	for i := range sequenceData{
		time := sequenceData[i].UpdatedAt
		fmt.Printf("%v>\n",time)
		if tmp.Before(time){
			t.Errorf("sequenceData error around %v",sequenceData[i].UpdatedAt)
		}
		tmp = time

	}
}