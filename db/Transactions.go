package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"time"
	"web_app/model"
)

//该事务检测用户是否存有文件
//若没有文件则将file_URL设为空

func AutoCheckEmptyURL() {
	{
		timer := time.Tick(time.Hour * 24)
		for {
			fmt.Println("-----start CheckEmptyURL------")
			err := CheckEmptyURL()
			if err != nil {
				fmt.Println("Error  when check empty URL: ", err)
			}
			<-timer
		}

	}
}

func CheckEmptyURL()error{
	db:=Link()
	defer db.Close()
	var datas []model.Demo_order

	db.Transaction(func(tx *gorm.DB) error {
		tx.Find(&datas)

		for i := range datas{
			path := "./tmp/"+datas[i].Order_no+"/file.txt"
			 _,err := os.Stat(path)
			 fmt.Println(err)
			 if err != nil{
				if os.IsNotExist(err){
					datas[i].File_url = ""
				}else{
					datas[i].File_url = "error"
				}
				 tx.Save(&datas[i])
			}
		}
		return nil
	})
return nil
}