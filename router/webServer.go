package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"web_app/handler"
	"web_app/model"
)

func WebServer(){
	fmt.Println("webServer is running")

	r := gin.Default()
	r.LoadHTMLFiles("templates/index.html")

	//测试
	r.GET("/hello", func(c *gin.Context) {

		//此处逻辑函数放handler
		c.JSON(200,gin.H{
			"message": 123,
		})
	})

	//创建数据
	r.POST("/json", func(c *gin.Context) {
		var dorder model.Demo_order
		if err := c.ShouldBindJSON(&dorder); err == nil{
			fmt.Printf("\nget data: %#v\n", dorder)

			handler.NewData(&dorder)
			c.JSON(http.StatusCreated, gin.H{
				"status":"get data",
			})
			}else{
			fmt.Println("get post failed! error:",err)
		}
	})

	//更新数据
	r.PUT("/json/:no", func(c *gin.Context) {
		var dorder model.Demo_order
		no := c.Param("no")

		fmt.Println("-----",no)



		if err := c.ShouldBindJSON(&dorder); err == nil && len(no) != 0{
			fmt.Printf("\nget data: %#v\n", dorder)

			if no != dorder.Order_no{
				fmt.Println(dorder.Order_no)
				fmt.Println("error order_NO has change")
				c.JSON(http.StatusBadRequest, gin.H{
					"status":"order_NO has change",
				})

				return
			}

			if err := handler.Update(no, &dorder); err ==nil{
				c.JSON(http.StatusCreated,gin.H{
					"status":no+"has change",
				})
			}else{
				c.JSON(http.StatusInternalServerError,gin.H{
					"status":"server bad!",
				})
			}

		}else{
			c.JSON(http.StatusUnauthorized,gin.H{
				"status":"error in data",
			})
		}

	})

	//获取某一数据
	r.GET("/json/:no", func(c *gin.Context) {
		no := c.Param("no")

		data := handler.SearchByNo(no)

		if data.IsEmpty(){
			c.JSON(http.StatusNotFound,gin.H{
				"status":"no such data",
			})
		}else{
			c.JSON(http.StatusOK,gin.H{
				"status":"GET",
				"order_no":data.Order_no,
				"user_name":data.User_name,
				"amount":data.Amount,
				"data_status":data.Status,
				"file_url":data.File_url,
				"create_time":data.CreatedAt,
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

		//按条件排序
		datas,_ := handler.Sequence(key)

		//模糊搜索name
		if search !=""{
			datas,_ = handler.Search(search,datas)
		}

		//格式化数据
		shows := make([]model.Demo_order_show,0,0)
		for _,data := range datas{
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

		os.MkdirAll(path, os.ModePerm)

		c.SaveUploadedFile(file,dst)

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

	r.Run(":8080")
}
