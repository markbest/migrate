package migrate

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "migrate/conf"
)

var db *gorm.DB

type Migrate struct {
	Id        int64  `gorm:"auto"`
	Migration string `gorm:"size(256)"`
	Batch     int64
}

func (Migrate) TableName() string {
	return Conf.Migrate.Table
}

func init() {
	db := DB()
	if !db.HasTable(&Migrate{}) {
		db.CreateTable(&Migrate{})
	}
}

func DB() *gorm.DB {
	if db == nil {
		newDb, err := newDB()
		if err != nil {
			panic(err)
		}
		newDb.DB().SetMaxIdleConns(10)
		newDb.DB().SetMaxOpenConns(100)
		newDb.LogMode(false)
		db = newDb
	}
	return db
}

func newDB() (*gorm.DB, error) {
	if err := InitConfig(); err != nil {
		fmt.Println(err)
	}

	sqlConnection := Conf.DB.User + ":" + Conf.DB.Password + "@tcp(" + Conf.DB.Host + ":" + Conf.DB.Port + ")/" + Conf.DB.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", sqlConnection)
	if err != nil {
		return nil, err
	}
	return db, nil
}

/* 获取所有已经执行的迁移文件 */
func GetAllMigrationsFile() (m []string) {
	var migrations []Migrate

	db = DB()
	db.Order("batch asc").Find(&migrations)
	if len(migrations) > 0 {
		for _, v := range migrations {
			m = append(m, v.Migration)
		}
	}
	return m
}

/* 获取最后一批操作的migrations */
func GetLatestMigrationsFile(action string) (batch int64, m []Migrate) {
	var migrate []Migrate

	db = DB()
	db.Order("batch desc").First(&migrate)

	if len(migrate) > 0 {
		for _, v := range migrate {
			batch = v.Batch
		}
	} else {
		batch = 0
	}

	if batch > 0 {
		if action == "up" {
			db.Where("batch <= ?", batch).Find(&m)
		} else if action == "down" {
			db.Where("batch = ?", batch).Find(&m)
		}
	}
	return batch, m
}

/* 执行迁移文件 */
func MigrateUp() {
	batch, m := GetLatestMigrationsFile("up")
	upsql, _, files := LoadMigrationsFile("up", m)

	db = DB()
	for _, k := range upsql {
		db.Exec(k)
	}

	if len(files) > 0 {
		for _, v := range files {
			db.Create(&Migrate{Migration: v, Batch: batch + 1})
			fmt.Print("migrate " + v + " successfully\n")
		}
	} else {
		fmt.Print("no migrations\n")
	}
}

/* 回滚迁移文件 */
func MigrateDown() {
	batch, m := GetLatestMigrationsFile("down")
	_, downsql, files := LoadMigrationsFile("down", m)

	db = DB()
	for _, k := range downsql {
		db.Exec(k)
	}

	if len(files) > 0 {
		for _, v := range files {
			db.Delete(&Migrate{Migration: v, Batch: batch})
			fmt.Print("rollback " + v + " successfully\n")
		}
	} else {
		fmt.Print("no rollback\n")
	}
}

func MigrateStatus() {
	//已经执行的迁移文件
	files := GetAllMigrationsFile()

	//尚未执行的迁移文件
	_, m := GetLatestMigrationsFile("up")
	_, _, others := LoadMigrationsFile("up", m)

	fmt.Println("+-----+----------------------------------------------------------+")
	fmt.Println("| Ran | Migration                                                |")
	fmt.Println("+-----+----------------------------------------------------------+")

	if len(files) > 0 {
		for _, v := range files {
			fmt.Println("|  Y  | " + v)
		}
	}

	if len(others) > 0 {
		for _, k := range others {
			fmt.Println("|  N  | " + k)
		}
	}

	fmt.Println("+-----+----------------------------------------------------------+")
}
