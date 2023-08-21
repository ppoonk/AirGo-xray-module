package model

import (
	"AirGo/global"
	"net/url"
	"strconv"
)

type Trojan struct {
	TrojanRemarks  string `json:"remarks"`
	TrojanPassword string `json:"password"`
	TrojanAddress  string `json:"address"`
	TrojanPort     int64  `json:"port"`

	TrojanNetwork  string `json:"network"` //传输协议 tcp,kcp,ws,h2,quic,grpc
	TrojanGrpcMode string `json:"mode"`    //grpc传输模式
	TrojanType     string `json:"type"`    //伪装类型 none http
	TrojanHost     string `json:"host"`    //伪装域名
	TrojanPath     string `json:"path"`    //

	TrojanSecurity      string `json:"security"` //传输层安全 none,tls,reality
	TrojanSni           string `json:"sni"`
	TrojanAlpn          string `json:"alpn"`
	TrojanAllowInsecure bool   `json:"allowInsecure"` //跳过证书验证

	TrojanFingerprint string `json:"fp"`
	TrojanPublicKey   string `json:"pbk"`
	TrojanShortId     string `json:"sid"`
	TrojanSpiderX     string `json:"spx"`
}

func ParseTrojanLink(link string) *Node {
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	if u.User == nil || u.Scheme != "trojan" {
		return nil
	}
	node := new(Node)
	node.NodeType = "trojan"
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
	if urlQuery.Get("network") != "" {
		node.Network = urlQuery.Get("network")
	}
	if urlQuery.Get("type") != "" {
		node.Type = urlQuery.Get("type")
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
	if urlQuery.Get("tls") != "" {
		node.Security = urlQuery.Get("tls")
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
