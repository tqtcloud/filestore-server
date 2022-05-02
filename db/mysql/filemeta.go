package mysqld

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

// Tab_file 初始SQL语句 实现更新字加载时间，软删除等
/** Tab_file
 create table `tbl_file` (
    `id` int(11) not null auto_increment ,
    `file_sha1` char(40) not null default '' comment '文件hash',
    `file_name` varchar(256) not null default '' comment '文件名',
    `file_size` bigint(20) default '0' comment '文件大小',
    `file_addr` varchar(1024) not null default '' comment '文件存储位置',
    `create_at` datetime DEFAULT CURRENT_TIMESTAMP  comment '创建日期',
    `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  comment '更新时间',
		`delete_at` datetime comment '更新时间',
    `status` int(11) not null default '0' comment '状态(可用/禁用/已删除)',
    `ext1` int(11) default '0' comment '备用字段1',
    `ext2` text comment '备用字段2' ,
    primary key (`id`) ,
    KEY `idx_status` (`status`),
    unique key `idx_file_hash`(`file_sha1`)
) engine=innodb default charset=utf8;
*/
type Tab_file struct {
	Id        int64          `gorm:"primaryKey;autoIncrement;type:int(11)"`
	File_sha1 string         `gorm:"column:file_sha1;type:char(40);not null;index;unique"`
	File_name string         `gorm:"column:file_name;type:varchar(256);not null;"`
	Fiel_size string         `gorm:"column:file_size;type:bigint(20)"`
	File_addr string         `gorm:"column:file_addr;not null;type:varchar(1024)"`
	CreateAT  *time.Time     `gorm:"column:create_at;autoCreateTime;default:null"`
	UpdateAT  *time.Time     `gorm:"column:update_at;autoCreateTime;default:null"`
	DeleteAt  gorm.DeletedAt `gorm:"column:delete_at"` //添加软删除
	Status    int            `gorm:"not null"`
}

// TableName 定义表名
func (f Tab_file) TableName() string {
	return "tbl_file"
}

var _db *gorm.DB

func init() {
	dsn := "root:123456@tcp(192.168.0.104:3306)/fileserver?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	sqlDB, _ := _db.DB()
	//设置数据库连接池参数
	sqlDB.SetConnMaxLifetime(100) //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)     //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。

}

// GetDB 获取gorm db对象，其他包需要执行数据库查询的时候，只要通过tools.getDB()获取db对象即可。
// GetDB 不用担心协程并发使用同样的db对象会共用同一个连接，db对象在调用他的方法的时候会从数据库连接池中获取新的连接
func (f Tab_file) GetDB() *gorm.DB {
	return _db
}

// FileMetaQuery 查找元数据返回结构体
func (f Tab_file) FileMetaQuery(meta *Tab_file) (*Tab_file, error) {
	db := f.GetDB()
	// 如果需要查找软删除数据请取消注释
	//err :=  db.Unscoped().Select("id,file_sha1,file_name,file_size,file_addr,create_at,update_at,status").Where("file_sha1 = ?",meta.File_sha1).Take(&meta).Error
	// 正常数据查找
	err := db.Select("id,file_sha1,file_name,file_size,file_addr,create_at,update_at,status").Where("file_sha1 = ?", meta.File_sha1).Take(&meta).Error
	// 当 First、Last、Take 方法找不到记录时，GORM 会返回 ErrRecordNotFound 错误
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//fmt.Println("查询不到数据",err)
		return meta, err
	} else if err != nil {
		//如果err不等于record not found错误，又不等于nil，那说明sql执行失败了。
		fmt.Println("查询失败", err)
		return meta, err
	}
	return meta, nil
}

// FileMetaInsert 插入数据
func (f Tab_file) FileMetaInsert(meta *Tab_file) error {
	db := f.GetDB()
	var err error
	_, err = f.FileMetaQuery(meta)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Create(&meta).Error; err != nil {
			fmt.Println("插入失败", err)
			return err
		}
		log.Printf("添加数据ID为：%d", meta.Id)
		return nil
	}
	return errors.New("插入数据存在，不再进行插入")
}

// FileMetaUpdate 更新元数据
func (f Tab_file) FileMetaUpdate(meta *Tab_file, oldUpdateName, newValue string) error {
	db := f.GetDB()
	var err error
	_, err = f.FileMetaQuery(meta)
	// 查询不到数据存在
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("数据不存在更新失败")
	}
	// 自定义键进行 Where 匹配
	if err = db.Model(&meta).Where("file_sha1 = ?", meta.File_sha1).Update(oldUpdateName, newValue).Error; err != nil {
		return err
	}
	log.Printf("元数据 %s 数据更新成功：%s -> %s", meta.File_sha1, oldUpdateName, newValue)
	return nil
}

func (f Tab_file) FileMetaDelete(meta *Tab_file) error {
	db := f.GetDB()
	err := db.Model(&meta).Where("file_sha1 = ?", meta.File_sha1).Delete(meta).Error
	if err != nil {
		return nil
	}
	return err
}
