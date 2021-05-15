package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/imagerecog"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)
var (
	file string
	filename string
	AccessKeyId string=""
	AccessKeySecret string=""
	Endpoint string=""
	Bucket string=""
	RegionId string=""
)

func httpsever(){
	server:=gin.Default()
	server.LoadHTMLGlob("static/*")
	server.POST("/result", func(c *gin.Context) {
	//获取表单数据 参数为name值
	f, err := c.FormFile("file")
	//错误处理
	if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{
	 "error": err,
	})
	return
	} else {
	//将文件保存至本项目根目录中

	c.SaveUploadedFile(f, "./img/"+f.Filename)
	file="./img/"+f.Filename
	upload()
	result:=detect()
	c.HTML(200,"result.html",gin.H{
		"src":file,
		"result":result,
	})
	}
	})
	server.GET("/", func(c *gin.Context) {
		c.HTML(200,"index.html","")
	})
	server.Run(":8080")
}

func detect() string{
	client, err := imagerecog.NewClientWithAccessKey(RegionId, AccessKeyId, AccessKeySecret)
	fileurl:="https://"+Bucket+"."+Endpoint+"/"+filename
	request := imagerecog.CreateClassifyingRubbishRequest()
	request.Scheme = "https"

	request.ImageURL = fileurl

	response, err := client.ClassifyingRubbish(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	info:=response.Data.Elements[0]
	return info.Category
}

func upload(){
	ossclient, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	tmp:=strings.Split(file,"/")
	filename=tmp[len(tmp)-1]
	// 获取存储空间。
	bucket, err := ossclient.Bucket(Bucket)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 上传本地文件。
	err = bucket.PutObjectFromFile(filename, file)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}

func main() {
	httpsever()
}
