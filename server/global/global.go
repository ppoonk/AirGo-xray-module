package global

import (
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Logrus *logrus.Logger
	//ClientWithSocks5 *http.Client
	//ClientWithDNS *http.Client
	//Dialer           *net.Dialer

	NodeAutoChangeCrontab *cron.Cron
	NodeAutoTcpingCrontab *cron.Cron

	Config ConfigInfo
)
