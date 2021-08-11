package handler

import (
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"time"
	"web_app/db"
	"web_app/model"
)

//创建数据
func NewData(data *model.Demo_order)error{
	fmt.Println("create new data")

	db := db.Link()
	defer db.Close()

	result := db.Create(data)
	if result.Error != nil{
		fmt.Println("Create error : ",result.Error)

		return result.Error

	}else{
		fmt.Println("Create success!")
	}

	return nil

}

//根据username搜索
func SearchByName(name string)(data *model.Demo_order){
	db := db.Link()
	defer db.Close()

	data.User_name = name

	db.First(data)

	return
}

//根据NO搜索
func SearchByNo(No string) *model.Demo_order{
	db := db.Link()
	defer db.Close()

	var data model.Demo_order

	data.Order_no = No

	db.First(&data)


	return &data
}

//更新数据
func Update(no string, newData *model.Demo_order)error{
	db := db.Link()
	defer db.Close()

	var data model.Demo_order

	data.Order_no = no
	db.First(&data)

	fmt.Println("-----",data)

	data.User_name = newData.User_name
	data.Amount = newData.Amount
	data.Status = newData.Status
	data.File_url = newData.File_url

	result :=db.Save(&data)
	if result.Error != nil{
		fmt.Println("Update error : ",result.Error)
		return result.Error
	}else{
		fmt.Println("Update success!")
	}

	return nil
}

//获取数据总数
func GetTotalNumber()int64{
	db := db.Link()
	defer db.Close()

	var datas model.Demo_order

	result := db.Find(&datas)

	totalNumber := result.RowsAffected

	return totalNumber
}

//获取指定长度数据
func GetData(start,end int)([]model.Demo_order, error){
	db := db.Link()
	defer db.Close()

	var datas []model.Demo_order

	result := db.Find(&datas)

	if result.Error != nil{
		fmt.Println("GetData error:",result.Error)
		return nil,result.Error
	}

	return datas, nil
}

//获取所有数据
func GetAllData()([]model.Demo_order, error){
	db := db.Link()
	defer db.Close()

	var datas []model.Demo_order

	result := db.Find(&datas)

	if result.Error != nil{
		fmt.Println("GetData error:",result.Error)
		return nil,result.Error
	}

	return datas, nil
}

//按amount排序
//从前往后遍历所有data，将amount项与后面所有项目比较
//找到后面amount的最大项与本项交换位置
//遍历到总数-1项，输出即为从大到小排序
func sequenceByAmount(datas []model.Demo_order){
	totalNum := GetTotalNumber()
	var i int64
	var t int64
	var Largest float64
	var target int64
	var tmp model.Demo_order

	for i=0;i<totalNum-1;i++{
		Largest = datas[i].Amount

		for t=i;t<totalNum;t++{
			if Largest < datas[t].Amount{
				Largest = datas[t].Amount
				target = t
			}
		}
		if i < target {		//该排序总是从前往后比较，防止t未赋值
			tmp = datas[i]
			datas[i] = datas[target]
			datas[target] = tmp
			//fmt.Println(i, "<=>", target)
		}
		target = i			//重置target
	}
}

//类似于以amount排序
func sequenceByTime(datas []model.Demo_order){
	totalNum := GetTotalNumber()
	var i int64
	var t int64
	var Largest time.Time
	var target int64
	var tmp model.Demo_order
	for i=0;i<totalNum-1;i++{
		Largest = datas[i].UpdatedAt
		for t=i;t<totalNum;t++{
			if Largest.Before(datas[t].UpdatedAt){
				Largest = datas[t].UpdatedAt
				target = t
			}
		}
		if i < target {
			tmp = datas[i]
			datas[i] = datas[target]
			datas[target] = tmp
			fmt.Println(i, "<=>", target)
			//fmt.Println(datas)
		}
		target = i
	}
}

//对order排序处理
func Sequence(key string)([]model.Demo_order, error){

	datas,_ := GetAllData()

	switch key {
	case "amount":
		sequenceByAmount(datas)
	case "time":
		sequenceByTime(datas)
	default:
		fmt.Println("key :",key)
	}

	return datas,nil
}

//按条件对name搜索
func Search(key string,datas []model.Demo_order)([]model.Demo_order, error){
	lenth := len(datas)

	//获取所有name并进行相似度排序
	names := make([]string,0,GetTotalNumber())
	for i :=range datas{
		names = append(names,datas[i].User_name)
	}
	resule := fuzzy.RankFind(key,names)
	sort.Sort(resule)

	//对排序好的结果写入新切片
	order_datas := make([]model.Demo_order,0,lenth)
	for i := range resule{
		order_datas = append(order_datas,datas[resule[i].OriginalIndex])
	}

	//furthest := func()int {
	//	t:=0
	//	for i := range resule{
	//		if t< resule[i].Distance{
	//			t = resule[i].Distance
	//		}
	//	}
	//	return t
	//}
	//
	//order_datas := make([]model.Demo_order,lenth,lenth)
	//
	//var t int64 = 0
	//for f:=furthest();f>=0;f--{
	//	for i:=range resule{
	//		if f == resule[i].Distance{
	//			order_datas[t] = datas[resule[i].OriginalIndex]
	//			t++
	//		}
	//	}
	//}

	//fmt.Println(resule)

	return order_datas,nil
}

//更新URL
func AddFileURL(no,URL string){
	order := SearchByNo(no)
	URL = URL+";\n"
	order.File_url = order.File_url + URL
	Update(no,order)
}