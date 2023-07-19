package model

import (
	"database/sql/driver"
	"encoding/json"
	"server/global"
	"time"
)

type Config struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DomesticType string `json:"domestic_type" gorm:"default:direct"` //国内分流类型，direct直连，proxy代理
	AbroadType   string `json:"abroad_type"   gorm:"default:direct"` //国外分流类型，direct直连，proxy代理

	Host               string `json:"host"`                                                //免流混淆
	AutoChangeNode     string `json:"auto_change_node" gorm:"default:0"`                   //自动切换节点,1启用
	WIFIProxy          string `json:"wifi_proxy"       gorm:"default:0;column:wifi_proxy"` //WiFi代理,1代理
	IPV6Net            string `json:"ipv6_net"         gorm:"default:1;column:ipv6_net"`   //ipv6联网，1联网
	AllowOutsideTcpUdp string `json:"allow_outside_tcp_udp" gorm:"default:1"`              //放行除tcp,udp外的流量,1放行
	AllowApps          Apps   `json:"allow_apps"`                                          //应用全局放行(填uid或包名)

}
type Apps []string

func (a *Apps) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, a)
}

func (a Apps) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// 获取配置,保存到全局变量
func (c *Config) GetConfig() error {
	var gc global.ConfigInfo
	err := global.DB.Model(&Config{}).First(&gc).Error
	if err != nil {
		return err
	}
	global.Config = gc
	err = global.DB.First(&c).Error
	if err != nil {
		return err
	}
	return nil
}

// 更新配置
func (c *Config) UpdateConfig() error {
	go func() {
		global.DB.Exec("update node_pool set host = ?", c.Host)
		global.DB.Exec("update node set host = ?", c.Host)
	}()
	err := global.DB.Save(&c).Error
	if err != nil {
		return err
	}
	//更新全局config
	var gc global.ConfigInfo
	err = global.DB.Model(&Config{}).First(&gc).Error
	if err != nil {
		return err
	}
	global.Config = gc
	return nil
}
