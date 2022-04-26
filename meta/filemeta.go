package meta

import (
	"fmt"
	"log"
)

// FileMeta 文件元数据结构体
type FileMeta struct {
	FileShah  string   //哈希
	FIleName string
	FileSize int64    //大小
	Location string   //上传地址
	UploadAt string   //时间
}

var fileMetas map[string]FileMeta

func init()  {
	fileMetas = make(map[string]FileMeta)
	log.Println("元数据初始化完毕...")
}

// UpdateFileMeta 新增/更新文件元数据
func UpdateFileMeta(fmeta FileMeta)  {
	// 哈希作为string 元信息作为数据传入
	fileMetas[fmeta.FileShah] = fmeta
	fmt.Println(fileMetas)
}

// GetFileMeta 通过shah值获取文件元数据对象
func GetFileMeta(fileshah string) FileMeta  {
	return fileMetas[fileshah]
}

// RemoveFileMeta 删除文件元数据
func RemoveFileMeta(fileSh1 string)  {
	delete(fileMetas,fileSh1)
	fmt.Println("删除文件后：" ,fileMetas)
}