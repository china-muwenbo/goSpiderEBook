# goSpiderEBook
golang 电子书网站爬虫  
golang 爬取电子书网站的数据库，并将数据存储在mongoDB中。 
电子书网站是：http://www.woaidu.org/ 
用到的技术有golang的协程池控制最大爬取协程数，如果不控制会因为资源不足而报错。  
电子书的数据的数据结构 

1.小说图片,存储url  
2.小说名字  
3.小说作者  
4.小说简介  
5.更新时间  
6.下载地址:(多个地址)                 
//数据在mongoDB的结构体                      
type EbookData struct {                   
  DataUrl  string                 
	ImageUrl string               
	Title    string                         
	Content  string               
	Time     string                     
	Author   string               
	Download string                 
}                                                       
//下载链接的结构体                                                 
type DownLoadInfo struct {    
	DownLoadUrl string  
	Downfrom    string  
	Updatetime  string  
	Progress    string  
}                                               
golang 连接mongoDB 的驱动用的是mgo: http://godoc.org/labix.org/v2/mgo 
  
