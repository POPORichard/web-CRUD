package handler

import (
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"time"
	"web_app/database"
	"web_app/model"
)

//NewData 创建数据
func NewData(data *model.DemoOrder)error{
	fmt.Println("create new data")

	db := database.Link()
	defer db.Close()

	result := db.Create(data)
	if result.Error != nil{
		fmt.Println("Create error : ",result.Error)

		return result.Error

	}

	fmt.Println("Create success!")

	return nil

}

//SearchByName 根据username搜索
func SearchByName(name string)(data *model.DemoOrder){
	db := database.Link()
	defer db.Close()

	data.UserName = name

	db.Where("user_name=?",name).First(data)

	return
}

//SearchByNo 根据NO搜索
func SearchByNo(No string) *model.DemoOrder {
	db := database.Link()
	defer db.Close()
	var data model.DemoOrder
	data.OrderNo = No
	db.Where("Order_no=?", No).First(&data)
	if data.ID == 0{
		return nil
	}
	fmt.Println("======",data)
	return &data
}

//Update 更新数据
func Update(no string, newData *model.DemoOrder)error{
	db := database.Link()
	defer db.Close()

	var data model.DemoOrder

	data.OrderNo = no
	db.Where("Order_no=?",no).First(&data)
	//此处应判断数据合法性
	data.UserName = newData.UserName
	data.Amount = newData.Amount
	data.Status = newData.Status
	data.FileURL = newData.FileURL


	result :=db.Save(&data)
	if result.Error != nil{
		fmt.Println("Update error : ",result.Error)
		return result.Error
	}

	fmt.Println("Update success!")

	return nil
}

//GetTotalNumber 获取数据总数
func GetTotalNumber()int64{
	db := database.Link()
	defer db.Close()

	var datas model.DemoOrder

	result := db.Find(&datas)

	totalNumber := result.RowsAffected

	return totalNumber
}

//CutData 截取分页
//若originData为nil则从数据库中获取数据
func CutData(page,pageSize int64, originData []model.DemoOrder)([]model.DemoOrder, error){

	if originData == nil{
		db := database.Link()
		defer db.Close()
		originData,_ = GetAllData()
	}

	largest := GetTotalNumber()
	start := page*pageSize
	end := (page+1)*pageSize
	if end > largest{
		end = largest
		start = largest - pageSize
	}

	data := originData[start:end]		//这里不能直接传回data

	returnData := make([]model.DemoOrder,0,pageSize)

	for i := range data{
		returnData = append(returnData,data[i])
	}

	return returnData, nil
}

//GetAllData 获取所有数据
func GetAllData()([]model.DemoOrder, error){
	db := database.Link()
	defer db.Close()

	var datas []model.DemoOrder

	result := db.Find(&datas)

	if result.Error != nil{
		fmt.Println("GetData error:",result.Error)
		return nil,result.Error
	}

	return datas, nil
}

//sequenceByAmount 按amount排序
//从前往后遍历所有data，将amount项与后面所有项目比较
//找到后面amount的最大项与本项交换位置
//遍历到总数-1项，输出即为从大到小排序
func sequenceByAmount(datas []model.DemoOrder){
	totalNum := len(datas)
	var i int
	var t int
	var Largest float64
	var target int
	var tmp model.DemoOrder

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

//sequenceByTime 按update time排序
func sequenceByTime(datas []model.DemoOrder){
	totalNum := len(datas)
	var i int
	var t int
	var Largest time.Time
	var target int
	var tmp model.DemoOrder
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

//Sequence 对order按关键项进行排序处理
func Sequence(key string,data []model.DemoOrder)([]model.DemoOrder, error){

	if data ==nil{
		data,_ = GetAllData()
	}

	switch key {
	case "amount":
		sequenceByAmount(data)
	case "time":
		sequenceByTime(data)
	default:
		fmt.Println("key :",key)
	}

	return data,nil
}

//Search 按条件对name模糊搜索
func Search(key string,datas []model.DemoOrder)([]model.DemoOrder, error){
	lenth := len(datas)

	//获取所有name并进行相似度排序
	names := make([]string,0,GetTotalNumber())
	for i :=range datas{
		names = append(names,datas[i].UserName)
	}
	resule := fuzzy.RankFind(key,names)
	sort.Sort(resule)

	//对排序好的结果写入新切片
	orderDatas := make([]model.DemoOrder,0,lenth)
	for i := range resule{
		orderDatas = append(orderDatas,datas[resule[i].OriginalIndex])
	}

	return orderDatas,nil
}

//AddFileURL 更新file_URL
func AddFileURL(no,URL string){
	order := SearchByNo(no)
	if order ==nil{
		return
	}
	order.FileURL = URL
	Update(no,order)
}