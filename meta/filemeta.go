package meta

import (
	mysqld "filestore1/db/mysql"
	"fmt"
	"log"
	"strconv"
)

// FileMeta 文件元数据结构体
type FileMeta struct {
	FileShah string //哈希
	FIleName string
	FileSize int64  //大小
	Location string //上传地址
	UploadAt string //时间
	//mysqld.Tab_file
}

// NewTabFile 将 FileMeta 中的数据 传递给 Tab_file
func NewTabFile(meta FileMeta) mysqld.Tab_file {
	return mysqld.Tab_file{
		File_sha1: meta.FileShah,
		File_name: meta.FIleName,
		Fiel_size: strconv.FormatInt(meta.FileSize, 10),
		File_addr: meta.Location,
	}
}

//var fileMetas = mysqld.Tab_file
var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
	log.Println("元数据初始化完毕...")
}

// UpdateFileMeta 新增/更新文件元数据
func UpdateFileMeta(fmeta FileMeta) {
	// 哈希作为string 元信息作为数据传入
	metatab := NewTabFile(fmeta)
	// 如果查询结果有错误则表明数据库中没有数据，则插入数据
	if _, err := metatab.FileMetaQuery(&metatab); err != nil {
		err := metatab.FileMetaInsert(&metatab)
		if err != nil {
			fmt.Printf("文件元数据插入数据库失败：%s\n", err)
			return
		}
		//fileMetas[fmeta.FileShah] = fmeta
		log.Println("文件元信息插入数据库成功")
		fmt.Println(metatab)
		return
	} else { // 否则为修改数据库内容
		err := metatab.FileMetaUpdate(&metatab, "file_name", fmeta.FIleName)
		if err != nil {
			fmt.Printf("元数据更新失败：%s\n", err)
		}
	}

}

// GetFileMeta 通过shah值获取文件元数据对象
func GetFileMeta(fileshah string) *mysqld.Tab_file {
	meta := FileMeta{}
	meta.FileShah = fileshah
	metatab := NewTabFile(meta)
	metatab1, err := metatab.FileMetaQuery(&metatab)
	if err != nil {
		fmt.Println(err)
		return &mysqld.Tab_file{}
	}
	return metatab1
}

// RemoveFileMeta 删除文件元数据
func RemoveFileMeta(fileSh1 string) {
	meta := mysqld.Tab_file{
		File_sha1: fileSh1,
	}
	delete(fileMetas, fileSh1)
	meta.FileMetaDelete(&meta)

	fmt.Println("删除文件后：", fileMetas)
}
