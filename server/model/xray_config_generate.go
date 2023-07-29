package model

import (
	"path/filepath"
	"runtime"
	"server/global"
	"server/utils"
	"strings"
)

// 生成xray-core配置文件
func GenerateConfig() error {
	var conf = map[string]interface{}{
		"log":       logConfig(),
		"inbounds":  inboundsConfig(),
		"outbounds": outboundConfig(),
		//"policy":    policyConfig(),
		"dns":     dnsConfig(),
		"routing": routingConfig(),
	}
	path := filepath.Join(global.Config.ExecutionPath, "config.json")
	err := utils.WriteJSON(conf, path)
	if err != nil {
		//panic(err)
		return err
	}
	//config, err := json.Marshal(&conf)
	//if err != nil {
	//	//global.Logrus.Error("生成xray-core配置文件错误:", err.Error())
	//	return err
	//}
	//json.RawMessage(config)
	//return string(config), nil
	return nil
}

// 路由
func routingConfig() interface{} {
	rules := make([]interface{}, 0)
	rules = append(rules, map[string]interface{}{
		"type":        "field",
		"network":     "udp",
		"port":        "53",
		"inboundTag":  []string{"udp-in"},
		"outboundTag": "dns-out",
	})
	rules = append(rules, map[string]interface{}{
		"type":        "field",
		"outboundTag": "direct",
		"domain": []string{
			"localhost",
		},
	})
	rules = append(rules, map[string]interface{}{
		"type": "field",
		"ip": []string{
			"geoip:private",
		},
		"outboundTag": "direct",
	})
	//国内分流
	switch global.Config.DomesticType {
	case "proxy":
		rules = append(rules, map[string]interface{}{
			"type": "field",
			"domain": []string{
				"geosite:cn",
				"geosite:apple",
				"geosite:microsoft",
			},
			"outboundTag": "domestic",
		})
		rules = append(rules, map[string]interface{}{
			"type": "field",
			"ip": []string{
				"geoip:cn",
			},
			"outboundTag": "domestic",
		})
	default: //设置出站直连
		rules = append(rules, map[string]interface{}{
			"type": "field",
			"domain": []string{
				"geosite:cn",
				"geosite:apple",
				"geosite:microsoft",
			},
			"outboundTag": "direct",
		})
		rules = append(rules, map[string]interface{}{
			"type": "field",
			"ip": []string{
				"geoip:cn",
			},
			"outboundTag": "direct",
		})
	}
	//国外分流
	//switch global.Config.AbroadType { //国外分流
	//case "proxy":
	//	rules = append(rules, map[string]interface{}{
	//		"type": "field",
	//		"domain": []string{
	//			"geosite:google",
	//			"geosite:facebook",
	//			"geosite:twitter",
	//			"geosite:telegram",
	//			"geosite:geolocation-!cn",
	//		},
	//		"outboundTag": "abroad",
	//	})
	//	rules = append(rules, map[string]interface{}{
	//		"type": "field",
	//		"ip": []string{
	//			"geoip:ae",
	//			"geoip:au",
	//			"geoip:br",
	//			"geoip:ca",
	//			"geoip:de",
	//			"geoip:dk",
	//			"geoip:es",
	//			"geoip:fi",
	//			"geoip:fr",
	//			"geoip:gb",
	//			"geoip:gr",
	//			"geoip:hk",
	//			"geoip:id",
	//			"geoip:il",
	//			"geoip:in",
	//			"geoip:iq",
	//			"geoip:ir",
	//			"geoip:it",
	//			"geoip:jp",
	//			"geoip:kr",
	//			"geoip:mo",
	//			"geoip:my",
	//			"geoip:nl",
	//			"geoip:no",
	//			"geoip:nz",
	//			"geoip:ph",
	//			"geoip:ru",
	//			"geoip:sa",
	//			"geoip:sg",
	//			"geoip:th",
	//			"geoip:tr",
	//			"geoip:tw",
	//			"geoip:us",
	//			"geoip:vn",
	//		},
	//		"outboundTag": "abroad",
	//	})
	//default: //设置出站直连
	//	rules = append(rules, map[string]interface{}{
	//		"type": "field",
	//		"domain": []string{
	//			"geosite:geolocation-!cn",
	//		},
	//		"outboundTag": "direct",
	//	})
	//}
	return map[string]interface{}{
		"domainStrategy": "IPIfNonMatch",
		"domainMatcher":  "hybrid",
		"rules":          rules,
	}
}

// 出站
func outboundConfig() interface{} {
	out := make([]interface{}, 0)
	//国外分流
	switch global.Config.AbroadType {
	case "proxy": //获取国外激活节点
		var node = &Node{Ascription: "abroad"}
		node, err := node.GetEnabledNodes()
		node.DomainToIP()
		if err != nil {
			return nil
		}
		out = append(out, NodeOutBound(node))
	default: //设置出站直连

	}
	//国内分流
	switch global.Config.DomesticType {
	case "proxy": //获取国内激活节点
		var node = &Node{Ascription: "domestic"}
		node, err := node.GetEnabledNodes()
		node.DomainToIP()
		if err != nil {
			return nil
		}
		out = append(out, NodeOutBound(node))
	default: //设置出站直连

	}

	out = append(out, map[string]interface{}{
		"tag":      "direct",
		"protocol": "freedom",
		"settings": map[string]interface{}{
			"domainStrategy": "UseIP",
			"userLevel":      0,
		},
	})
	out = append(out, map[string]interface{}{
		"tag":      "block",
		"protocol": "blackhole",
		"settings": map[string]interface{}{
			"response": map[string]interface{}{
				//"type": "http",
				"type": "none",
			},
		},
	})
	out = append(out, map[string]interface{}{
		"tag":      "dns-out",
		"protocol": "dns",
	})
	return out
}

// 根据节点类型生成出站配置
func NodeOutBound(n *Node) interface{} {
	switch n.NodeType {
	case "trojan":
		return trojanOutbound(*n)
	case "vmess":
		return vMessOutbound(*n)
	default:
		return vLessOutbound(*n)
	}
}

// Trojan
func trojanOutbound(trojan Node) interface{} {
	streamSettings := map[string]interface{}{
		"network":  "tcp",
		"security": "tls",
	}
	streamSettings["tlsSettings"] = map[string]interface{}{
		"allowInsecure": trojan.AllowInsecure,
		"serverName":    trojan.Sni,
	}
	return map[string]interface{}{
		"tag":      trojan.Ascription,
		"protocol": "trojan",
		"settings": map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address":  trojan.Ns,
					"port":     trojan.Port,
					"password": trojan.UUID,
					"level":    0,
				},
			},
		},
		"streamSettings": streamSettings,
	}
}

// VMess
func vMessOutbound(n Node) interface{} {
	streamSettings := map[string]interface{}{
		"network":  n.Network,
		"security": n.Security,
	}
	if n.Security == "tls" {
		tlsSettings := map[string]interface{}{
			"allowInsecure": false,
		}
		if n.Sni != "" {
			tlsSettings["serverName"] = n.Sni
		}
		if n.Alpn != "" {
			tlsSettings["alpn"] = strings.Split(n.Alpn, ",")
		}
		streamSettings["tlsSettings"] = tlsSettings
	}
	switch n.Network {
	case "tcp":
		streamSettings["tcpSettings"] = map[string]interface{}{
			"header": map[string]interface{}{
				"type": n.Type,
			},
		}
	case "kcp":
		kcpSettings := map[string]interface{}{
			"mtu":              1350,
			"tti":              50,
			"uplinkCapacity":   12,
			"downlinkCapacity": 100,
			"congestion":       false,
			"readBufferSize":   2,
			"writeBufferSize":  2,
			"header": map[string]interface{}{
				"type": n.Type,
			},
		}
		if n.Type != "none" {
			kcpSettings["seed"] = n.Path
		}
		streamSettings["kcpSettings"] = kcpSettings
	case "ws":
		streamSettings["wsSettings"] = map[string]interface{}{
			"path": n.Path,
			"headers": map[string]interface{}{
				"Host": n.Host,
			},
		}
	case "h2":
		host := make([]string, 0)
		for _, line := range strings.Split(n.Host, ",") {
			line = strings.TrimSpace(line)
			if line != "" {
				host = append(host, line)
			}
		}
		streamSettings["httpSettings"] = map[string]interface{}{
			"path": n.Path,
			"host": host,
		}
	case "quic":
		quicSettings := map[string]interface{}{
			"security": n.Host,
			"header": map[string]interface{}{
				"type": n.Type,
			},
		}
		if n.Host != "none" {
			quicSettings["key"] = n.Path
		}
		streamSettings["quicSettings"] = quicSettings
	case "grpc":
		streamSettings["grpcSettings"] = map[string]interface{}{
			"serviceName": n.Path,
			"multiMode":   n.Type == "multi",
		}
	}
	return map[string]interface{}{
		"tag":      n.Ascription,
		"protocol": "vmess",
		"settings": map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": n.Ns,
					"port":    n.Port,
					"users": []interface{}{
						map[string]interface{}{
							"id":       n.UUID,
							"alterId":  n.Aid,
							"security": n.Scy,
							"level":    0,
						},
					},
				},
			},
		},
		"streamSettings": streamSettings,
		"mux": map[string]interface{}{
			"enabled": false,
		},
	}
}

// VLESS
func vLessOutbound(vless Node) interface{} {
	//mux := setting.Mux()
	mux := false
	security := vless.Security
	network := vless.Network
	user := map[string]interface{}{
		"id":         vless.UUID,
		"flow":       vless.VlessFlow,
		"encryption": vless.VlessEncryption,
		"level":      0,
	}
	streamSettings := map[string]interface{}{
		"network":  network,
		"security": security,
	}
	switch security {
	case "tls":
		tlsSettings := map[string]interface{}{
			"allowInsecure": vless.AllowInsecure,
		}
		sni := vless.Sni
		alpn := vless.Alpn
		if sni != "" {
			tlsSettings["serverName"] = sni
		}
		if alpn != "" {
			tlsSettings["alpn"] = strings.Split(alpn, ",")
		}
		streamSettings["tlsSettings"] = tlsSettings
	case "xtls":
		xtlsSettings := map[string]interface{}{
			"allowInsecure": false,
		}
		sni := vless.Sni
		alpn := vless.Alpn
		if sni != "" {
			xtlsSettings["serverName"] = sni
		}
		if alpn != "" {
			xtlsSettings["alpn"] = strings.Split(alpn, ",")
		}
		streamSettings["xtlsSettings"] = xtlsSettings
		mux = false
	case "reality":
		realitySettings := map[string]interface{}{
			"show": false,
			//"fingerprint": vless.GetValue(field.FingerPrint),
			//"serverName":  vless.GetHostValue(field.SNI),
			//"publicKey":   vless.GetValue(field.PublicKey),
			//"shortId":     vless.GetValue(field.ShortId),
			//"spiderX":     vless.GetValue(field.SpiderX),
		}
		streamSettings["realitySettings"] = realitySettings
		mux = false
	}
	switch network {
	case "tcp":
		streamSettings["tcpSettings"] = map[string]interface{}{
			"header": map[string]interface{}{
				"type": vless.Type,
			},
		}
	case "kcp":
		kcpSettings := map[string]interface{}{
			"mtu":              1350,
			"tti":              50,
			"uplinkCapacity":   12,
			"downlinkCapacity": 100,
			"congestion":       false,
			"readBufferSize":   2,
			"writeBufferSize":  2,
			"header": map[string]interface{}{
				"type": vless.Type, //伪装类型
			},
		}
		if vless.Type != "none" {
			kcpSettings["seed"] = vless.Path
		}
		streamSettings["kcpSettings"] = kcpSettings
	case "h2":
		mux = false
		host := make([]string, 0)
		for _, line := range strings.Split(vless.Host, ",") {
			line = strings.TrimSpace(line)
			if line != "" {
				host = append(host, line)
			}
		}
		streamSettings["httpSettings"] = map[string]interface{}{
			"path": vless.Path,
			"host": host,
		}
	case "ws":
		streamSettings["wsSettings"] = map[string]interface{}{
			"path": vless.Path,
			"headers": map[string]interface{}{
				"Host": vless.Host,
			},
		}
	case "quic":
		quicSettings := map[string]interface{}{
			"security": vless.Security,
			"header": map[string]interface{}{
				"type": vless.Type,
			},
		}
		if vless.Host != "none" {
			quicSettings["key"] = vless.Path
		}
		streamSettings["quicSettings"] = quicSettings
	case "grpc":
		streamSettings["grpcSettings"] = map[string]interface{}{
			"serviceName": vless.Path,
			"multiMode":   "multi",
		}
	}
	return map[string]interface{}{
		"tag":      vless.Ascription,
		"protocol": "vless",
		"settings": map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": vless.Ns,
					"port":    vless.Port,
					"users": []interface{}{
						user,
					},
				},
			},
		},
		"streamSettings": streamSettings,
		"mux": map[string]interface{}{
			"enabled": mux,
		},
	}
}

// 日志
func logConfig() interface{} {
	//path := filepath.Join(utils.GetConfigDir(), "xray_access.log")
	return map[string]string{
		//"access":   path,
		"loglevel": "warning",
	}
}

// 入站
func inboundsConfig() interface{} {
	data := []interface{}{
		map[string]interface{}{
			"tag":      "redir-tcp",
			"port":     1230, //!
			"protocol": "dokodemo-door",
			"sniffing": map[string]interface{}{
				"enabled": true, //!
				"destOverride": []string{
					"http",
					"tls",
				},
			},
			"settings": map[string]interface{}{
				"network":        "tcp,udp",
				"followRedirect": true,
			},
		},
	}
	switch runtime.GOOS {
	case "linux":
		//data = append(data, map[string]interface{}{
		//	"tag":      "udp-in",
		//	"port":     1231,
		//	"protocol": "dokodemo-door",
		//	"settings": map[string]interface{}{
		//		"auth": "noauth",
		//		"udp":  true,
		//	},
		//	"sniffing": map[string]interface{}{
		//		"enabled": true,
		//		"destOverride": []string{
		//			"http",
		//			"tls",
		//		},
		//	},
		//})
		data = append(data, map[string]interface{}{
			"tag":      "udp-in",
			"port":     1231,
			"protocol": "dokodemo-door",
			"settings": map[string]interface{}{
				"network":        "udp",
				"followRedirect": true,
			},
			"sniffing": map[string]interface{}{
				"enabled":      true,
				"destOverride": []string{"http", "tls"},
			},
			"streamSettings": map[string]interface{}{
				"sockopt": map[string]interface{}{
					"tproxy": "tproxy",
				},
			},
		})
	default:
		data = append(data, map[string]interface{}{
			"tag":      "udp-in",
			"port":     1231,
			"protocol": "socks",
			"settings": map[string]interface{}{
				"auth": "noauth",
				"udp":  true,
			},
			"sniffing": map[string]interface{}{
				"enabled": true,
				"destOverride": []string{
					"http",
					"tls",
				},
			},
		})
	}
	return data
}

// 本地策略
func policyConfig() interface{} {
	return map[string]interface{}{
		"levels": map[string]interface{}{
			"0": map[string]interface{}{
				"handshake":    4,
				"connIdle":     300,
				"uplinkOnly":   1,
				"downlinkOnly": 1,
				"bufferSize":   10240,
			},
		},
		"system": map[string]interface{}{
			"statsInboundUplink":   true,
			"statsInboundDownlink": true,
		},
	}
}

// DNS
func dnsConfig() interface{} {
	servers := make([]interface{}, 0)
	servers = append(servers, "223.6.6.6")
	return map[string]interface{}{
		"servers": servers,
	}
}
