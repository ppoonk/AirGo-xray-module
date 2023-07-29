package initialize

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"path/filepath"
	"server/global"
)

// 初始化sqlite数据库
func InitGormSqlite() *gorm.DB {

	path := filepath.Join(global.Config.ExecutionPath, "airgo.db")
	global.Logrus.Info("数据库路径:", path)
	if db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		SkipDefaultTransaction: true, //关闭事务，将获得大约 30%+ 性能提升
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "gormv2_",
			SingularTable: true, //单数表名
		},
	}); err != nil {
		//global.Logrus.Error("gorm.Open error:", err)
		panic(err)
	} else {
		//sqlDB, _ := db.DB()
		//sqlDB.SetMaxIdleConns(global.Config.Mysql.MaxIdleConns)
		//sqlDB.SetMaxOpenConns(global.Config.Mysql.MaxOpenConns)
		return db
	}
}
