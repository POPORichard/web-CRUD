package router

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"web_app/database"
	"web_app/handler"
	"web_app/model"
)

var URL string = "/json"

var reader = strings.NewReader(`{"order_no":"test","user_name":"test","amount":888,"status":"test"}`)

func deleteTestData(order *model.DemoOrder) {
	db := database.Link()

	if order == nil{
		order = handler.SearchByNo("test")
	}
	db.Unscoped().Delete(order)
	db.Close()
}

func TestWebServer(t *testing.T) {

	r := WebServer()
	w :=httptest.NewRecorder()

	req,err :=http.NewRequest("POST","/json",reader)
	if err != nil{
		t.Fatal("Request err:",err)
	}

	r.ServeHTTP(w,req)

	assert.Equal(t,200,w.Code)

	time.Sleep(time.Second *3)






	//if err != nil{
	//	t.Fatal("Create test data failed! err:",err)
	//}else if resp.StatusCode != http.StatusOK{
	//	t.Fatal("Response code is", resp.StatusCode)
	//}else{
	//	data := handler.SearchByNo("test")
	//	if data ==nil{
	//		t.Fatal("data did not write into database!")
	//	}else if data.User_name != "test" || data.Amount != 888 || data.Status != "test"{
	//		t.Fatal("data has changed!")
	//	}else{
	//		fmt.Println("data write in success!")
	//		defer deleteTestData(data)
	//	}
	//}

}
