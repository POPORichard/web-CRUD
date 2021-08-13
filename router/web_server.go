package router

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"web_app/handler"
	"web_app/model"
	"web_app/server"
)
// WebServer 网络服务
func WebServer() *gin.Engine{
	fmt.Println("webServer is running")

	r := gin.Default()
	r.LoadHTMLFiles("templates/index.html")

	//测试
	r.GET("/hello", func(c *gin.Context) {

		c.JSON(200,gin.H{
			"message": 123,
		})
	})

	//创建数据
	r.POST("/json", func(c *gin.Context) {
		var order model.DemoOrder
		if err := c.ShouldBindJSON(&order); err == nil{
			fmt.Printf("\nget data: %#v\n", order)

			handler.NewData(&order)
			c.JSON(http.StatusCreated, gin.H{
				"status":"get data",
			})
			}else{
			fmt.Println("get post failed! error:",err)
		}
	})

	//更新数据
	r.PUT("/json/:no", func(c *gin.Context) {

		no := c.Param("no")
		data := handler.SearchByNo(no)

		b,err := c.GetRawData()  // 从c.Request.Body读取请求数据
		if err != nil {
			fmt.Println("Error can't get data! err:",err)
			return
		}
		var m map[string]interface{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			fmt.Println("Error can't get data! err:",err)
			return
		}

		orderNo,possess := m["orderNo"];if !possess{
			fmt.Println("Can't get No from body")
			c.JSON(http.StatusUnauthorized,gin.H{
				"status":"error in data",
			})
			return
		}

		if orderNo != data.OrderNo {
			fmt.Println("error order_NO is different")
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "order_NO is different with request",
			})
			return
		}

		userName, possess := m["userName"]
		if possess {
			data.UserName = userName.(string)
		}

		amount, possess := m["amount"]
		if possess {
			data.Amount = amount.(float64)
		}

		status, possess := m["status"]
		if possess {
			data.Status = status.(string)
		}


		if err := handler.Update(no, data); err == nil {
			c.JSON(http.StatusCreated, gin.H{
				"status": no + "has change",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "server bad!",
			})
		}

		//if err := c.ShouldBindJSON(&order); err == nil && len(no) != 0{
		//	fmt.Printf("\nget data: %#v\n", order)
		//
		//	if no != order.OrderNo {
		//		fmt.Println(order.OrderNo)
		//		fmt.Println("error order_NO has change")
		//		c.JSON(http.StatusBadRequest, gin.H{
		//			"status":"order_NO has change",
		//		})
		//		return
		//	}
		//
		//	if err := handler.Update(no, &order); err ==nil{
		//		c.JSON(http.StatusCreated,gin.H{
		//			"status":no+"has change",
		//		})
		//	}else{
		//		c.JSON(http.StatusInternalServerError,gin.H{
		//			"status":"server bad!",
		//		})
		//	}
		//
		//}else{
		//	c.JSON(http.StatusUnauthorized,gin.H{
		//		"status":"error in data",
		//	})
		//}
	})

	//获取某一数据
	r.GET("/json/:no", func(c *gin.Context) {
		no := c.Param("no")

		fmt.Println("===查询数据=== No:",no)

		data := handler.SearchByNo(no)

		if data.IsEmpty(){
			c.JSON(http.StatusNotFound,gin.H{
				"status":"no such data",
			})
		}else{
			c.JSON(http.StatusOK,gin.H{
				"status":"GET",
				"order_no":data.OrderNo,
				"user_name":data.UserName,
				"amount":data.Amount,
				"data_status":data.Status,
				"file_url":data.FileURL,
				"create_time":data.CreatedAt,
				"update_time":data.UpdatedAt,
			})
		}
	})

	//获取数据列表
	r.GET("/index", func(c *gin.Context) {

		//分页
		//start := c.DefaultQuery("start","0")
		//end := c.Query("end")
		//lim := c.Query("lim")
		//page := c.DefaultQuery("page","1")
		key := c.Query("key")
		search := c.DefaultQuery("search","")

		allData,err := handler.GetAllData()
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"status":"error! can not get data!",
			})
		}

		//按条件排序
		allData,_ = handler.Sequence(key, allData)

		//模糊搜索name
		if search !=""{
			allData,_ = handler.Search(search, allData)
		}

		//格式化数据
		shows := make([]model.DemoOrderShow,0,0)
		for _,data := range allData {
			shows = append(shows,data.OrderToShow())
		}

		c.HTML(http.StatusOK, "index.html",shows)

	})

	//文件上传
	r.POST("/upload/:no", func(c *gin.Context) {
		no := c.Param("no")
		file, err:= c.FormFile("f1")
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"status":err.Error(),
			})
			return
		}

		log.Println(file.Filename)

		path := "./tmp/"+no+"/"
		dst := fmt.Sprintf(path+"file.txt")

		err = os.MkdirAll(path, os.ModePerm)
		if err != nil{
			c.JSON(http.StatusServiceUnavailable,gin.H{
				"status":err.Error(),
			})
			return
		}

		err = c.SaveUploadedFile(file,dst)
		if err != nil{
			c.JSON(http.StatusServiceUnavailable,gin.H{
				"status":err.Error(),
			})
			return
		}

		handler.AddFileURL(no,"http://127.0.0.1:8080/download/"+no+"/file.txt")

		c.JSON(http.StatusOK, gin.H{
			"status":fmt.Sprintf("'%s' uploaded success!", file.Filename),
		})
	})

	//下载文件
	r.GET("/download/:no/:filename", func(c *gin.Context) {
		no := c.Param("no")
		filename := c.Param("filename")
		path := "./tmp/"+no+"/"+filename


		c.Header("Content-Type","application/txt")
		c.Header("Content-Disposition","attachment; filename=\"" + filename + "\"")
		c.File(path)

	})

	//下载excel
	r.GET("/list", func(c *gin.Context) {
		var wg sync.WaitGroup
		wg.Add(1)

		err := server.WriteToExcel(nil, &wg)
		if err != nil{
			c.JSON(http.StatusServiceUnavailable,gin.H{
				"status":"unable to create excel err:"+err.Error(),
			})
		}

		time.Sleep(time.Second)

		path := "../list.xlsx"
		c.Header("Content-Type","application/txt")
		c.Header("Content-Disposition","attachment; filename=list.xlsx")
		c.File(path)


	})

	return r
}
