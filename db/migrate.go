package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "github.com/markbest/migrate/conf"
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

/* 获取最后一批操作的迁移 */
func GetLatestMigrationsFile(action string) (batch int64, m []Migrate) {
	var migrate Migrate
	db = DB()
	db.Order("batch desc").First(&migrate)

	if migrate.Batch > 0 {
		batch = migrate.Batch
	} else {
		batch = 0
	}

	if batch > 0 {
		if action == "up" {
			db.Where("batch <= ?", batch).Find(&m)
		}
		if action == "down" {
			db.Where("batch = ?", batch).Find(&m)
		}
	}
	return batch, m
}

/* 执行迁移 */
func HandleMigrateUp() {
	batch, m := GetLatestMigrationsFile("up")
	upSql, _, files := LoadMigrationsFile("up", m)

	db = DB()
	tx := db.Begin()
	for k, v := range upSql {
		if err := tx.Exec(v).Error; err != nil {
			tx.Rollback()
			panic(err)
		} else {
			db.Create(&Migrate{Migration: files[k], Batch: batch + 1})
			fmt.Println("migrate " + files[k] + " successfully")
		}
	}
	tx.Commit()

	if len(files) == 0 {
		fmt.Println("no migration files")
	}
}

/* 回滚迁移 */
func HandleMigrateDown() {
	_, m := GetLatestMigrationsFile("down")
	_, downSql, files := LoadMigrationsFile("down", m)

	db = DB()
	tx := db.Begin()
	for k, v := range downSql {
		if err := tx.Exec(v).Error; err != nil {
			tx.Rollback()
			panic(err)
		} else {
			db.Where("migration = ?", files[k]).Delete(Migrate{})
			fmt.Println("rollback " + files[k] + " successfully")
		}
	}
	tx.Commit()

	if len(files) == 0 {
		fmt.Println("no rollback files")
	}
}

func HandleMigrateStatus() {
	files := GetAllMigrationsFile()
	_, m := GetLatestMigrationsFile("up")
	_, _, others := LoadMigrationsFile("up", m)

	fmt.Println("+-----+----------------------------------------------------------+")
	fmt.Println("| Ran | Migration                                                |")
	fmt.Println("+-----+----------------------------------------------------------+")

	if len(files) > 0 {
		for _, v := range files {
			fmt.Println("|  Y  | " + v)
		}
	} else {
		fmt.Println("| No migrate files |")
	}

	if len(others) > 0 {
		for _, k := range others {
			fmt.Println("|  N  | " + k)
		}
	}

	fmt.Println("+-----+----------------------------------------------------------+")
}