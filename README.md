数据库连接语句在db.link中

功能说明：
创建数据
    /json
    POST 
    Header Content-Type：application/json
    Body exp：{
                "order_no":"017",
                "user_name":"BBBCC",
                "amount":8,
                "status":"green"
               }
                
更新数据
    /json/:order_no
    PUT
    Header Content-Type：application/json
    Body exp：{
                "order_no":"017",
                "user_name":"BBBCC",
                "amount":8,
                "status":"green"
               }
                    
查询数据
    /json/:order_no
    GET
    
上传文件
    /upload/:order_no
    POST
    Header Content-Type multipart/form-data
    Body    f1:file.txt
    
列表
    /index?key=amount&search=***
    key:amount 按amount排序
        time    按更新时间排序
    search:模糊查找user_name
    
列表中file_URL为文件下载地址
    
    
    

