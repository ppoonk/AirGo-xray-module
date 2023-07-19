package model

import (
	"encoding/json"
	"fmt"
	"server/global"
	"server/utils"
	"strconv"
	"strings"
)

type Vmess struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port int    `json:"port"`
	Id   string `json:"id"`
	Scy  string `json:"scy"`  //加密方式 auto,none,chacha20-poly1305,aes-128-gcm,zero
	Aid  int    `json:"aid"`  //额外ID
	Net  string `json:"net"`  //传输协议 tcp,kcp,ws,h2,quic,grpc
	Type string `json:"type"` //伪装类型 none,http
	Host string `json:"host"`
	Path string `json:"path"`

	Tls  string `json:"tls"`
	Sni  string `json:"sni"`
	Alpn string `json:"alpn"`
}

func ParseVMessLink(link string) *Node {
	node := new(Node)
	node.NodeType = "vmess"
	if strings.ToLower(link[:8]) == "vmess://" {
		link = link[8:]
	} else {
		return nil
	}
	if len(link) == 0 {
		return nil
	}
	jsonStr := utils.SubBase64Decode(link)
	if jsonStr == "" {
		return nil
	}
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	if err != nil {
		return nil
	}
	if version, ok := mapResult["v"]; ok {
		node.V = fmt.Sprintf("%v", version)
	}
	if ps, ok := mapResult["ps"]; ok {
		node.Remarks = fmt.Sprintf("%v", ps) //别名
	} else {
		return nil
	}
	if addr, ok := mapResult["add"]; ok {
		node.Address = fmt.Sprintf("%v", addr) //地址
	} else {
		return nil
	}
	if scy, ok := mapResult["scy"]; ok {
		node.Scy = fmt.Sprintf("%v", scy) //加密方式 auto,none,chacha20-poly1305,aes-128-gcm,zero
	} else {
		node.Scy = "auto"
	}
	if port, ok := mapResult["port"]; ok {
		value, err := strconv.Atoi(fmt.Sprintf("%v", port))
		if err == nil {
			node.Port = value //端口
		} else {
			return nil
		}
	} else {
		return nil
	}

	if id, ok := mapResult["id"]; ok {
		node.UUID = fmt.Sprintf("%v", id) //uuid
	} else {
		return nil
	}
	if aid, ok := mapResult["aid"]; ok {
		if value, err := strconv.Atoi(fmt.Sprintf("%v", aid)); err == nil {
			node.Aid = value //额外id
		} else {
			return nil
		}
	} else {
		return nil
	}
	if net, ok := mapResult["net"]; ok {
		node.Network = fmt.Sprintf("%v", net) //传输协议
	} else {
		return nil
	}
	if type1, ok := mapResult["type"]; ok {
		node.Type = fmt.Sprintf("%v", type1)
	} else {
		return nil
	}

	//获取混淆
	if global.Config.Host != "" {
		node.Host = global.Config.Host
	} else if host, ok := mapResult["host"]; ok {
		node.Host = fmt.Sprintf("%v", host)
	} else {
		return nil
	}

	if path, ok := mapResult["path"]; ok {
		node.Path = fmt.Sprintf("%v", path)
	} else {
		return nil
	}
	if tls, ok := mapResult["tls"]; ok {
		node.Security = fmt.Sprintf("%v", tls)
	} else {
		return nil
	}
	if sni, ok := mapResult["sni"]; ok {
		node.Sni = fmt.Sprintf("%v", sni)
	}
	if alpn, ok := mapResult["alpn"]; ok {
		node.Alpn = fmt.Sprintf("%v", alpn)
	}
	return node
}
