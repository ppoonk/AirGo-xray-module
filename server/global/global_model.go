package global

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type ConfigInfo struct {
	ID        int64     `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	OS            string `json:"os"`                                  //系统os linux android darwin
	ExecutionPath string `json:"execution_path"`                      //执行文件路径
	DomesticType  string `json:"domestic_type" gorm:"default:direct"` //国内分流类型，direct直连，proxy代理
	AbroadType    string `json:"abroad_type"   gorm:"default:direct"` //国外分流类型，direct直连，proxy代理

	Host               string `json:"host"`                                                //免流混淆
	StartupXray        string `json:"startup_xray" gorm:"default:0"`                       //开机启动xray，1启动
	BlockAds           string `json:"block_ads"        gorm:"default:0"`                   //广告域名屏蔽
	NodePoolModel      string `json:"node_pool_model"  gorm:"default:hm"`                  //节点池工作模式，手动（hm）自动（am）负载均衡（bm）
	WIFIProxy          string `json:"wifi_proxy"       gorm:"default:1;column:wifi_proxy"` //WiFi代理,1代理
	IPV6Net            string `json:"ipv6_net"         gorm:"default:0;column:ipv6_net"`   //ipv6联网，1联网
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
