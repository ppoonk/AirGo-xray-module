package initialize

import (
	"runtime"
	"server/global"
	"server/model"
	"server/utils"
)

// 注册数据库表
func RegisterTables() {
	err := global.DB.AutoMigrate(
		model.Node{},
		model.NodePool{},
		model.Subscription{},
		model.Config{},
		global.ConfigInfo{},
	)
	if err != nil {
		//os.Exit(0)
		global.Logrus.Error("table AutoMigrate error:", err.Error())
		return
	}
	global.Logrus.Info("table AutoMigrate success")
}
func InsertInto() {
	if global.DB.Migrator().HasTable(&model.Subscription{}) {
		if err := global.DB.First(&model.Subscription{}).Error; err == nil {
			return
		}
	}
	//创建系统订阅
	subData := model.Subscription{Alias: "系统订阅"}
	err := global.DB.Create(&subData).Error
	if err != nil {
		global.Logrus.Error(err.Error())
		return
	}
	//创建配置
	configData := model.Config{}
	switch runtime.GOOS {
	case "darwin":
		configData.OS = "darwin"
	default:
		var sh model.Shell
		out, _ := sh.DoShell(model.GetAndroidVersion, true)
		if out != "" {
			configData.OS = "android"
		} else {
			configData.OS = "linux"
		}
	}
	configData.ExecutionPath = utils.GetRunPath()
	err = global.DB.Create(&configData).Error
	if err != nil {
		global.Logrus.Error(err.Error())
		return
	}
	//globalConfigData := global.ConfigInfo{}
	//err = global.DB.Create(&globalConfigData).Error
	//if err != nil {
	//	global.Logrus.Error(err.Error())
	//	return
	//}
}
