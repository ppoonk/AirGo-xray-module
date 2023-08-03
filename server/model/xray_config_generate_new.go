package model

import "server/global"

// 出站
func outboundConfigNew() interface{} {
	out := make([]interface{}, 0)
	//国外分流
	switch global.Config.AbroadType {
	case "proxy": //获取国外激活节点
		switch global.Config.NodePoolModel {
		case "bm": //负载均衡
			var node = &Node{Ascription: "abroad"}
			nodeArr, _ := node.GetNodePool()
			for _, v := range *nodeArr {
				v.DomainToIP()
				out = append(out, NodeOutBound(&v))
			}
		default: //其他模式
			var node = &Node{Ascription: "abroad"}
			node, err := node.GetEnabledNodes()
			node.DomainToIP()
			if err != nil {
				return nil
			}
			out = append(out, NodeOutBound(node))
		}

	default: //设置出站直连

	}
	//国内分流
	switch global.Config.DomesticType {
	case "proxy": //获取国内激活节点
		switch global.Config.NodePoolModel {
		case "bm": //负载均衡
			var node = &Node{Ascription: "domestic"}
			nodeArr, _ := node.GetNodePool()
			for _, v := range *nodeArr {
				v.DomainToIP()
				out = append(out, NodeOutBound(&v))
			}
		default: //其他模式
			var node = &Node{Ascription: "domestic"}
			node, err := node.GetEnabledNodes()
			node.DomainToIP()
			if err != nil {
				return nil
			}
			out = append(out, NodeOutBound(node))
		}

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

// 路由
func routingConfigNew() interface{} {
	//处理负载均衡tag
	balancers := make([]interface{}, 0)
	switch global.Config.NodePoolModel {
	case "bm":
		var domesticTags, abroadTags []string
		var node1 = &Node{Ascription: "domestic"}
		var node2 = &Node{Ascription: "abroad"}
		nodeArr1, _ := node1.GetNodePool()
		for _, v := range *nodeArr1 {
			domesticTags = append(domesticTags, v.Remarks)
		}
		nodeArr2, _ := node2.GetNodePool()
		for _, v := range *nodeArr2 {
			abroadTags = append(abroadTags, v.Remarks)
		}
		balancers = append(balancers, map[string]interface{}{
			"tag":      "domestic",
			"selector": domesticTags,
		})
		balancers = append(balancers, map[string]interface{}{
			"tag":      "abroad",
			"selector": abroadTags,
		})
	default:
	}
	//处理私网和dns
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
		switch global.Config.NodePoolModel {
		case "bm":
			rules = append(rules, map[string]interface{}{
				"type": "field",
				"domain": []string{
					"geosite:cn",
					"geosite:apple",
					"geosite:microsoft",
				},
				"balancerTag": "domestic",
			})
			rules = append(rules, map[string]interface{}{
				"type": "field",
				"ip": []string{
					"geoip:cn",
				},
				"balancerTag": "domestic",
			})
		default:
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
		}

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
	//国外分流，负载均衡时使用，其他情况默认出站为国外
	switch global.Config.AbroadType {
	case "proxy":
		switch global.Config.NodePoolModel {
		case "bm":
			rules = append(rules, map[string]interface{}{
				"type": "field",
				"domain": []string{
					"geosite:google",
					"geosite:facebook",
					"geosite:twitter",
					"geosite:telegram",
					"geosite:geolocation-!cn",
				},
				"balancerTag": "abroad",
			})
			rules = append(rules, map[string]interface{}{
				"type": "field",
				"ip": []string{
					"geoip:ae",
					"geoip:au",
					"geoip:br",
					"geoip:ca",
					"geoip:de",
					"geoip:dk",
					"geoip:es",
					"geoip:fi",
					"geoip:fr",
					"geoip:gb",
					"geoip:gr",
					"geoip:hk",
					"geoip:id",
					"geoip:il",
					"geoip:in",
					"geoip:iq",
					"geoip:ir",
					"geoip:it",
					"geoip:jp",
					"geoip:kr",
					"geoip:mo",
					"geoip:my",
					"geoip:nl",
					"geoip:no",
					"geoip:nz",
					"geoip:ph",
					"geoip:ru",
					"geoip:sa",
					"geoip:sg",
					"geoip:th",
					"geoip:tr",
					"geoip:tw",
					"geoip:us",
					"geoip:vn",
				},
				"balancerTag": "abroad",
			})
		default:

		}
	default:

	}
	switch global.Config.NodePoolModel {
	case "bm":
		return map[string]interface{}{
			"domainStrategy": "IPIfNonMatch",
			"domainMatcher":  "hybrid",
			"balancers":      balancers,
			"rules":          rules,
		}
	default:
		return map[string]interface{}{
			"domainStrategy": "IPIfNonMatch",
			"domainMatcher":  "hybrid",
			"rules":          rules,
		}
	}
}
