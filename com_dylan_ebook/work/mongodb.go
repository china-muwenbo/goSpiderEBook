package work



import (
	"fmt"
	"gopkg.in/mgo.v2"
)


const mongodburl = "mongodb://127.0.0.1:27017"
const Dbname="ebook"
const Document="ebook"

type Operater struct {
	Mogsession *mgo.Session
	Dbname     string
	Document   string
}

//连接数据库
func (operater *Operater) Connect() error {
	mgo.SetDebug(true)
	//mgo.SetLogger()
	mgo.SetStats(true)
	mogsession, err := mgo.Dial(mongodburl)
	if err != nil {
		fmt.Println(err)
		return err
	}
	operater.Mogsession =mogsession

	return nil
}

//插入
func (operater *Operater) Insert(doc interface{}) error {
	collcetion:=operater.Mogsession.DB(operater.Dbname).C(operater.Document)

	err := collcetion.Insert(doc)
	return err
}

//插入
func (operater *Operater) Close() {
	operater.Mogsession.Close()
}

//统计文档中数据的个数
func (operater *Operater) Count() (int,error) {
	collcetion:=operater.Mogsession.DB(operater.Dbname).C(operater.Document)
	i,err:=collcetion.Count()
	return i,err
}


