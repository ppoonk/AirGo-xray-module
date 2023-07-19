package global

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Logrus *logrus.Logger
	//ClientWithSocks5 *http.Client
	//ClientWithDNS *http.Client
	//Dialer           *net.Dialer

	Config ConfigInfo
)
