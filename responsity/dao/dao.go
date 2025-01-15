package dao

import (
	"MiniPrograms/responsity/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MiniProgramsDAO struct {
	db *gorm.DB
}

// 初始化 SQLite 数据库连接
func InitDB() (*gorm.DB, error) {
	// 使用 SQLite 连接数据库，这里可以指定 SQLite 数据库文件的路径
	db, err := gorm.Open(sqlite.Open("./db/mini_programs.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&model.MiniPrograms{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 创建一个新的 MiniPrograms 实例
func NewMiniProgramsDAO(db *gorm.DB) *MiniProgramsDAO {
	return &MiniProgramsDAO{db: db}
}

func (d *MiniProgramsDAO) Find(name string) (*model.MiniPrograms, error) {
	var miniPrograms model.MiniPrograms
	err := d.db.Model(&model.MiniPrograms{}).First(&miniPrograms).Where("name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &miniPrograms, nil
}

func (d *MiniProgramsDAO) Save(miniPrograms model.MiniPrograms) error {
	return d.db.Model(&model.MiniPrograms{}).Where("name= ?", miniPrograms.Name).Save(&miniPrograms).Error
}
