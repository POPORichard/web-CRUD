package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"web_app/database"
	"web_app/handler"
	"web_app/model"
)

var URL string = "/json"

var reader = strings.NewReader(`{"orderNo":"test","userName":"test","amount":888,"status":"test"}`)
//var readerChange = strings.NewReader(`{"orderNo":"test","userName":"test","amount":111,"status":"test"}`)

func deleteTestData(order *model.DemoOrder) {
	db := database.Link()

	if order == nil{
		order = handler.SearchByNo("test")
	}
	if order == nil{
		panic("LooK out! Going to delete database!")
	}
	db.Unscoped().Delete(order)
	db.Close()
}

func TestWebServer(t *testing.T) {

	r := WebServer()
	w :=httptest.NewRecorder()

	req,err :=http.NewRequest("POST","/json",reader)
	if err != nil{
		t.Fatal("Make new request err:",err)
	}

	r.ServeHTTP(w,req)

	if w == nil{
		t.Fatal("Create test data failed! Nothing in recorder!")
	}else if w.Code != 201{
		t.Fatal("Create data failed! Status code is:",w.Code)
	}

	data := handler.SearchByNo("test")
	if data ==nil{
		t.Fatal("data did not write into database!")
	}else if data.UserName != "test" || data.Amount != 888 || data.Status != "test"{
				t.Fatal("data has changed!")
	}
	fmt.Println("data write in success!")

	defer deleteTestData(data)

	//数据更新测试
	req,err = http.NewRequest("PUT","/json/test",strings.NewReader(`{"orderNo":"test","amount":111}`))
	if err != nil{
		t.Fatal("Make new request err:",err)
	}

	w =httptest.NewRecorder()
	r.ServeHTTP(w,req)
	if w == nil{
		t.Fatal("Update data failed! Nothing in recorder!")
	}else if w.Code != 201{
		t.Fatal("Update data failed! Status code is:",w.Code)
	}

	data = handler.SearchByNo("test")
	if data.UserName != "test" || data.Amount != 111 || data.Status != "test"{
		fmt.Println(data)
		t.Fatal("data has changed!")
	}
	fmt.Println("Data update is success!")

	//数据查询测试
	req,err = http.NewRequest("GET","/json/test",nil)
	if err != nil{
		t.Fatal("Make new request err:",err)
	}

	w =httptest.NewRecorder()
	r.ServeHTTP(w,req)
	if w == nil{
		t.Fatal("Search test data failed! Nothing in recorder!")
	}else if w.Code != 200{
		t.Fatal("Search data failed! Status code is:",w.Code)
	}

	if data.UserName != "test" || data.Amount != 111 || data.Status != "test"{
		t.Fatal("data has changed!")
	}
	fmt.Println("Data update in success!")







}
