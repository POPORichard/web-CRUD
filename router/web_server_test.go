package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
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

	//写入数据
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

	res := w.Result()
	respond,_ := ioutil.ReadAll(res.Body)
	var respondData model.DemoOrderShow
	err = json.Unmarshal(respond,&respondData)
	if err != nil{
		t.Fatal("Get data not expect err:",err)
	}

	if !reflect.DeepEqual(respondData,data.OrderToShow()){
		fmt.Println(data)
		fmt.Println(respondData)
		t.Fatal("Search data is not equal with data in database!")
	}

	//文件上传测试
	path := "../../tmp/test.txt"
	file, err := os.Open(path)
	if err != nil {
		t.Fatal("Open file err:",err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("f1", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	req,err = http.NewRequest("POST","/upload/test",body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w =httptest.NewRecorder()
	r.ServeHTTP(w,req)
	if w == nil{
		t.Fatal("Upload file failed! Nothing in recorder!")
	}else if w.Code != 201{
		t.Fatal("Upload file failed! Status code is:",w.Code)
	}

	data = handler.SearchByNo("test")
	if data.FileURL != "http://127.0.0.1:8080/download/test/file.txt"{
		t.Fatal("data.File_URL is not right! Which is :",data.FileURL)
	}
	fmt.Println("Data update is success!")

	//文件下载测试
	req,err = http.NewRequest("GET","/download/test/file.txt",nil)
	w =httptest.NewRecorder()
	r.ServeHTTP(w,req)
	if w == nil{
		t.Fatal("Download file failed! Nothing in recorder!")
	}else if w.Code != 200{
		t.Fatal("Download file failed! Status code is:",w.Code)
	}

	_,err = os.Stat("./tmp/test/file/txt")
	if err != nil{
		os.IsExist(err)
		fmt.Println("Download file is exist!")
	}else{
		t.Fatal("Download file is not exist")
	}

	err = os.RemoveAll("./tmp")
	if err != nil{
		t.Error("Delete tmp file failed!")
	}


	//生成及下载Excel测试
	req,err = http.NewRequest("GET","/list",nil)
	w =httptest.NewRecorder()
	r.ServeHTTP(w,req)
	if w == nil{
		t.Fatal("Upload file failed! Nothing in recorder!")
	}else if w.Code != 200{
		t.Fatal("Upload file failed! Status code is:",w.Code)
	}


}
