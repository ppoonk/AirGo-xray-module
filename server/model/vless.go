package model

import (
	"AirGo/global"
	"net/url"
	"strconv"
)

type Vless struct {
	VlessRemarks string `json:"remarks"` //别名
	VlessID      string `json:"id"`      //用户id
	VlessAddress string `json:"address"` //地址
	VlessPort    int64  `json:"port"`    //端口

	VlessFlow       string `json:"flow"`       //流控 none,xtls-rprx-vision,xtls-rprx-vision-udp443
	VlessEncryption string `json:"encryption"` //加密方式 none

	VlessNetwork  string `json:"network"` //传输协议 tcp,kcp,ws,h2,quic,grpc
	VlessGrpcMode string `json:"mode"`    //grpc传输模式
	VlessType     string `json:"type"`    //伪装类型 none http
	VlessHost     string `json:"host"`    //伪装域名
	VlessPath     string `json:"path"`    //path

	VlessSecurity      string `json:"security"` //传输层安全 none,tls,reality
	VlessSni           string `json:"sni"`
	VlessAlpn          string `json:"alpn"`
	VlessAllowInsecure bool   `json:"allowInsecure"` //跳过证书验证

	VlessFingerprint string `json:"fp"`
	VlessPublicKey   string `json:"pbk"`
	VlessShortId     string `json:"sid"`
	VlessSpiderX     string `json:"spx"`
}

// vless://uuid@abc:80?encryption=none&type=ws&security=&host=www.10086.com&path=%2Fpath#%E5%B1%B1%E4%B8%9C
// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
func ParseVLessLink(link string) *Node {
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	if u.User == nil || u.Scheme != "vless" {
		return nil
	}
	node := new(Node)
	node.NodeType = "vless"

	//remarks
	node.Remarks = u.Fragment
	if node.Remarks == "" {
		node.Remarks = u.Host
	}
	//address
	node.Address = u.Hostname()
	//port
	node.Port, err = strconv.ParseInt(u.Port(), 10, 64)
	if err != nil {
		return nil
	}
	//uuid
	node.UUID = u.User.Username()

	//解析参数
	urlQuery := u.Query()
	if urlQuery.Get("flow") != "" {
		node.VlessFlow = urlQuery.Get("flow")
	}
	if urlQuery.Get("encryption") != "" {
		node.VlessEncryption = urlQuery.Get("encryption")
	}
	if urlQuery.Get("type") != "" {
		node.Network = urlQuery.Get("type")
	}
	if urlQuery.Get("security") != "" {
		node.Security = urlQuery.Get("security")
	}
	//获取混淆
	if global.Config.Host != "" {
		node.Host = global.Config.Host
	} else if urlQuery.Get("host") != "" {
		node.Host = urlQuery.Get("host")
	} else {
		return nil
	}

	if urlQuery.Get("path") != "" {
		node.Path = urlQuery.Get("path")
	}

	if urlQuery.Get("sni") != "" {
		node.Sni = urlQuery.Get("sni")
	}
	if urlQuery.Get("alpn") != "" {
		node.Alpn = urlQuery.Get("alpn")
	}
	if urlQuery.Get("allowInsecure") != "" {
		node.AllowInsecure = true
	}
	return node
}
