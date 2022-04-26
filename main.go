package main

import (
	"filestore1/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Server start ...")

	addr := "127.0.0.1:7070"
	http.HandleFunc("/file/upload",handlers.UploadHandler)
	http.HandleFunc("/file/upload/suc",handlers.UploadSucHandler)
	http.HandleFunc("/file/meta",handlers.GetFileMeteHandler)
	http.HandleFunc("/file/download",handlers.DownloadHandler)
	http.HandleFunc("/file/update",handlers.FileMetaUpdateHandier)
	http.HandleFunc("/file/delete",handlers.FileDeleteHandler)
	// 开启端口
	err := http.ListenAndServe(addr, nil)
	if err != nil{
		panic(fmt.Errorf("端口冲突：%s\n",err))
	}
	log.Println("Server start :7070")
}
