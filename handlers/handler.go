package handlers

import (
	"encoding/json"
	"filestore1/meta"
	"filestore1/util"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request )  {
	dir ,_ := os.Getwd()
	if r.Method == "GET"{
		//返回上传html页面
		data,err := ioutil.ReadFile("./static/view/index.html")
		if err != nil{
			io.WriteString(w,"internel server ell")
			return
		}
		io.WriteString(w,string(data))
	} else if r.Method == "POST" {
		// 接受文件流及存储到本地目录
		file,head,err := r.FormFile("file")
		if err != nil{
			fmt.Printf("Failed to get data,err: %s\n",err.Error())
			return
		}
		defer file.Close()

		fileMete := meta.FileMeta{
			FileShah: head.Filename,
			Location: fmt.Sprintf("%s",filepath.Join(dir,"tmp",head.Filename)),  //拼接全局路径
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newfile,err := os.Create(fileMete.Location)
		if err != nil{
			fmt.Printf("Failed to create file err: %s\n",err.Error())
		}
		defer newfile.Close()

		// 返回文件大小和错误
		fileMete.FileSize,err = io.Copy(newfile,file)
		if err != nil{
			fmt.Printf("Failed to save into file,err: %s\n",err.Error())
			return
		}

		//计算哈希值
		newfile.Seek(0,0)
		fileMete.FileShah = util.FileSha1(newfile)
		meta.UpdateFileMeta(fileMete)

		http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
	}
}

// UploadSucHandler 文件上传完成跳转页面
func UploadSucHandler(w http.ResponseWriter,r *http.Request)  {
	io.WriteString(w,"文件上传成功！！")
}

// GetFileMeteHandler 获取文件元信息
func GetFileMeteHandler(w http.ResponseWriter , r *http.Request)  {
	// 解析form
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fMate := meta.GetFileMeta(filehash)
	data,err := json.Marshal(fMate)
	if err != nil{
		fmt.Printf("Meta json Marshal err: %s\n ",err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter , r *http.Request)  {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	// 读取元数据中的文件指定路径，如果没有则报错
	f,err := os.Open(fm.Location)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data , err := ioutil.ReadAll(f)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//设置http响应头，这样才能浏览器识别下载
	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("Content-Descrption","attachment;filename=\""+fm.FIleName+"\"")
	w.Write(data)
}

// FileMetaUpdateHandier 文件重命名
func FileMetaUpdateHandier(w http.ResponseWriter ,r *http.Request)  {
	r.ParseForm()
	//  opType 0 为更名操作，其他值为其他操作
	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST"{  //必须大写
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	// 更新元数据的文件名称时，连带路径中的也更新了  暂时还未能做到
	//dir,_:=os.Getwd()
	curFileMeta.FIleName = newFileName
	//os.FileInfo()
	//os.Rename()
	//curFileMeta.Location = filepath.Join(dir,"tmp",newFileName)

	meta.UpdateFileMeta(curFileMeta) //更新元数据

	data,err := json.Marshal(curFileMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
// FileDeleteHandler 删除文件
func FileDeleteHandler(w http.ResponseWriter , r *http.Request)  {
	r.ParseForm()
	filesh1 := r.Form.Get("filehash")

	fMeta := meta.GetFileMeta(filesh1)
	fmt.Println(fMeta.Location)

	err := os.Remove(fMeta.Location)
	if err == nil{
		log.Printf("文件已删除：%s",fMeta.Location)
	}else {
		fmt.Println(err)
	}
	meta.RemoveFileMeta(filesh1)


	w.WriteHeader(http.StatusOK)

}