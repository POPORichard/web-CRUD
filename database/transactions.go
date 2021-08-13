package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"time"
	"web_app/model"
)

//AutoCheckEmptyURL 该事务检测用户是否存有文件
//若没有文件则将file_URL设为空
func AutoCheckEmptyURL() {
	{
		timer := time.Tick(time.Hour * 24)
		for {
			fmt.Println("-----start CheckEmptyURL------")
			err := checkEmptyURL()
			if err != nil {
				fmt.Println("Error  when check empty URL: ", err)
			}
			<-timer
		}

	}
}

func checkEmptyURL()error{
	db:=Link()
	defer db.Close()
	var datas []model.DemoOrder

	db.Transaction(func(tx *gorm.DB) error {
		tx.Find(&datas)

		for i := range datas{
			path := "./tmp/"+datas[i].OrderNo +"/file.txt"
			 _,err := os.Stat(path)
			 fmt.Println(err)
			 if err != nil{
				if os.IsNotExist(err){
					datas[i].FileURL = ""
				}else{
					datas[i].FileURL = "error"
				}
				 tx.Save(&datas[i])
			}
		}
		return nil
	})
return nil
}