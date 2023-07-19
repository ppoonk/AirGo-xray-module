package model

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"server/global"
	"server/utils"
	"strconv"
	"strings"
	"time"
)

const (
	ipTest1 = "https://myip.ipip.net"
	ipTest2 = "http://3.0.3.0/ip"
)

type Node struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt gorm.DeletedAt `json:"-"  gorm:"index"`

	SubscriptionID uint   `json:"subscription_id"` //foreign key
	NodeType       string `json:"node_type"`       //节点类型

	//基础参数
	Remarks    string  `json:"remarks"`    //别名
	UUID       string  `json:"uuid"`       //用户id
	Address    string  `json:"address"`    //地址
	Port       int     `json:"port"`       //端口
	Ns         string  `json:"ns"`         //ip地址
	TcpingData float64 `json:"tcping"`     //
	Ascription string  `json:"ascription"` //abroad domestic
	Enabled    bool    `json:"enabled"`    //是否为激活节点

	//vmess参数
	V   string `json:"v"`
	Scy string `json:"scy"` //加密方式 auto,none,chacha20-poly1305,aes-128-gcm,zero
	Aid int    `json:"aid"` //额外ID
	//vless参数
	VlessFlow       string `json:"flow"`       //流控 none,xtls-rprx-vision,xtls-rprx-vision-udp443
	VlessEncryption string `json:"encryption"` //加密方式 none

	//传输参数
	Network  string `json:"network"` //传输协议 tcp,kcp,ws,h2,quic,grpc
	Type     string `json:"type"`    //伪装类型 ws,h2：无    tcp,kcp：none，http    quic：none，srtp，utp，wechat-video，dtls，wireguard
	Host     string `json:"host"`    //伪装域名
	Path     string `json:"path"`    //path
	GrpcMode string `json:"mode"`    //grpc传输模式 gun，multi

	//传输层安全
	Security      string `json:"security"`                          //传输层安全类型 none,tls,reality
	Sni           string `json:"sni"`                               //
	Fingerprint   string `json:"fp"`                                //
	Alpn          string `json:"alpn"`                              //tls
	AllowInsecure bool   `json:"allowInsecure" gorm:"default:true"` //tls 跳过证书验证

	PublicKey string `json:"pbk"` //reality
	ShortId   string `json:"sid"` //reality
	SpiderX   string `json:"spx"` //reality
}
type NodePool struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt gorm.DeletedAt `json:"-"  gorm:"index"`

	SubscriptionID uint   `json:"subscription_id"` //foreign key
	NodeType       string `json:"node_type"`       //节点类型

	//基础参数
	Remarks    string  `json:"remarks"`    //别名
	UUID       string  `json:"uuid"`       //用户id
	Address    string  `json:"address"`    //地址
	Port       int     `json:"port"`       //端口
	Ns         string  `json:"ns"`         //ip地址
	TcpingData float64 `json:"tcping"`     //
	Ascription string  `json:"ascription"` //abroad domestic
	Enabled    bool    `json:"enabled"`    //是否为激活节点

	//vmess参数
	V   string `json:"v"`
	Scy string `json:"scy"` //加密方式 auto,none,chacha20-poly1305,aes-128-gcm,zero
	Aid int    `json:"aid"` //额外ID
	//vless参数
	VlessFlow       string `json:"flow"`       //流控 none,xtls-rprx-vision,xtls-rprx-vision-udp443
	VlessEncryption string `json:"encryption"` //加密方式 none

	//传输参数
	Network  string `json:"network"` //传输协议 tcp,kcp,ws,h2,quic,grpc
	Type     string `json:"type"`    //伪装类型 ws,h2：无    tcp,kcp：none，http    quic：none，srtp，utp，wechat-video，dtls，wireguard
	Host     string `json:"host"`    //伪装域名
	Path     string `json:"path"`    //path
	GrpcMode string `json:"mode"`    //grpc传输模式 gun，multi

	//传输层安全
	Security      string `json:"security"`                          //传输层安全类型 none,tls,reality
	Sni           string `json:"sni"`                               //
	Fingerprint   string `json:"fp"`                                //
	Alpn          string `json:"alpn"`                              //tls
	AllowInsecure bool   `json:"allowInsecure" gorm:"default:true"` //tls 跳过证书验证

	PublicKey string `json:"pbk"` //reality
	ShortId   string `json:"sid"` //reality
	SpiderX   string `json:"spx"` //reality
}

// 解析一条节点,vmess vless trojan
func ParseLink(link string) *Node {
	//fmt.Println("解析一条链接", link)
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	switch u.Scheme {
	case "vmess":
		if obj := ParseVMessLink(link); obj != nil {
			return obj
		}
	case "vless":
		if obj := ParseVLessLink(link); obj != nil {
			return obj
		}
	case "trojan":
		if obj := ParseTrojanLink(link); obj != nil {
			return obj
		}
	}
	return nil
}

// 根据id 查节点
func (n *Node) FindNodeById() error {
	return global.DB.First(&n).Error
}

// 节点域名解析成ip
func (n *Node) DomainToIP() {
	address := net.ParseIP(n.Address)
	if address != nil {
		// 是IP
		//fmt.Println("判断ip:", address)
		n.Ns = n.Address
	} else {
		// 自定义解析ip
		ipArr, err := utils.Resolver("223.6.6.6", 5*time.Second).LookupHost(context.Background(), n.Address)
		if err != nil {
			global.Logrus.Error("节点域名解析成ip,error:", err.Error())
			return
		}
		if len(ipArr) >= 1 {
			n.Ns = ipArr[0]
		} else {
			global.Logrus.Error("节点域名解析成ip,error")
			n.Ns = n.Address
			return
		}
	}
}

// 删除节点
func (n *Node) DeleteNode() error {
	return global.DB.Debug().Delete(&n, n.ID).Error
}

// 修改节点,保存到节点池
func (n *NodePool) UpdateNode() error {
	return global.DB.Save(&n).Error
}

// 修改节点,保存到节点池
func (n *Node) UpdateNode() error {
	return global.DB.Save(&n).Error
}

// 新建节点
func (n *Node) NewNode() error {
	return global.DB.Create(&n).Error
}

// 节点tcp测试
func (n *Node) Tcping() float64 {
	count := 0
	var sum float64

	var duraArr []int64
	for i := 0; i < 5; i++ {
		//d := net.Dialer{Timeout: 5 * time.Second}
		start := time.Now()
		//conn, err := global.Dialer.Dial("tcp", fmt.Sprintf("%s:%d", n.Address, n.Port))
		conn, err := utils.Dialer("223.6.6.6", 5*time.Second).Dial("tcp", fmt.Sprintf("%s:%d", n.Address, n.Port))
		if err != nil {
			continue
		}
		end := time.Since(start)
		conn.Close()
		count += 1
		duraArr = append(duraArr, int64(end))
	}
	//fmt.Println("次数：", count)
	//fmt.Println("时间间隔数组：", duraArr)
	for _, v := range duraArr {
		sum = sum + float64(v)
	}
	switch count {
	case 0:
		n.TcpingData = -1
		go n.UpdateNode()
		return -1
	default:
		n.TcpingData, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", sum/1e6/float64(count)), 64)
		go n.UpdateNode()
		return n.TcpingData
	}
	//
	//if count == 0 {
	//	n.TcpingData = -1
	//	go n.UpdateNode()
	//	return -1
	//} else if count == 1 {
	//	sum = float64(duraArr[0])
	//	n.TcpingData, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", sum/1e6/float64(count-2)), 64)
	//	go n.UpdateNode()
	//	return n.TcpingData
	//} else {
	//	r := utils.QuickSort(duraArr, 0, len(duraArr))
	//	for _, v := range r[:len(r)-1] {
	//		sum = sum + float64(v)
	//	}
	//	n.TcpingData, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", sum/1e6/float64(count-1)), 64)
	//	go n.UpdateNode()
	//	return n.TcpingData
	//}
}

// 保存激活的节点到数据库
func (n *Node) SetEnableNodeToDb() error {
	err := global.DB.Debug().Exec("UPDATE node_pool SET enabled = NULL WHERE enabled = ? and ascription = ?", n.Enabled, n.Ascription).Error
	if err != nil {
		return err
	}
	err = global.DB.Debug().Exec("UPDATE node_pool SET enabled = ? WHERE id = ?", n.Enabled, n.ID).Error
	if err != nil {
		return err
	}
	return nil
}

// 获取国内(国外)节点池
func (n *Node) GetNodePool() (*[]NodePool, error) {
	var nodePool []NodePool
	err := global.DB.Where(&NodePool{Ascription: n.Ascription}).Find(&nodePool).Error
	return &nodePool, err
}

// 加入节点池
func (n *Node) JoinNodePool() error {
	fmt.Println("加入节点池", n)
	n.ID = 0
	return global.DB.Debug().Model(&NodePool{}).Create(&n).Error
}

// 从节点池删除节点
func (n *Node) DeleteNodePool() error {
	return global.DB.Model(&NodePool{}).Delete(&n, n.ID).Error
}

// 设置激活的节点 ,判断当前激活的节点
// 如果有，能连通则跳过，否则遍历除本节点之外的节点池，设置延迟最低的一个
// 如果没有，则遍历节点池，设置延迟最低的一个
func (n *Node) SetEnableNode() bool {
	var node, enabledNode Node
	var nodes []Node
	err := global.DB.Model(&NodePool{}).Where("enabled = '1' and ascription =?", n.Ascription).First(&node).Error
	if err != nil {
		err1 := global.DB.Model(&NodePool{}).Where("ascription = ?", n.Ascription).Find(&nodes).Error
		if err1 != nil || len(nodes) == 0 {
			return false
		} else {
			//遍历节点池
			for _, v := range nodes {
				v.Tcping()
				if v.TcpingData > 0 {
					if enabledNode.ID != 0 {
						if v.TcpingData < enabledNode.TcpingData {
							enabledNode = v
							continue
						}
						continue
					}
					enabledNode = v
					continue
				}
				continue
			}
			//保存激活的节点到数据库
			if enabledNode.ID != 0 {
				global.Logrus.Info("节点切换为：", enabledNode.Remarks)
				enabledNode.Enabled = true
				enabledNode.SetEnableNodeToDb()
				return true
			} else {
				return false
			}
		}
	} else {
		//httping 判断连通性
		var d int64
		switch n.Ascription {
		case "abroad":
			d = n.TestNodeDelayNew("https://www.youtube.com")
		default:
			d = n.TestNodeDelayNew("https://www.baidu.com")
		}

		if d != -1 {
			global.Logrus.Info("连通,跳过切换", n.Ascription)
			return false
		} else {
			err2 := global.DB.Model(&NodePool{}).Where("ascription =?", n.Ascription).Find(&nodes).Error
			if err2 != nil || len(nodes) == 0 {
				return false
			} else {
				//遍历节点池
				for _, v := range nodes {
					v.Tcping()
					if v.TcpingData > 0 {
						if enabledNode.ID != 0 {
							if v.TcpingData < enabledNode.TcpingData {
								enabledNode = v
								continue
							}
							continue
						}
						enabledNode = v
						continue
					}
					continue
				}
				//保存激活的节点到数据库
				if enabledNode.ID != 0 {
					global.Logrus.Info("节点切换为：", enabledNode.Remarks)
					enabledNode.Enabled = true
					enabledNode.SetEnableNodeToDb()
					return true
				} else {
					return false
				}
			}
		}
	}
}

// 获取激活的国内国外节点
func (n *Node) GetEnabledNodes() (*Node, error) {

	var node Node
	var err error
	var errText string
	if n.Ascription == "domestic" {
		errText = "国内节点未配置"
		err = global.DB.Model(&NodePool{}).Where("enabled = '1' and ascription ='domestic'").First(&node).Error
	} else {
		errText = "国外节点未配置"
		err = global.DB.Model(&NodePool{}).Where("enabled = '1' and ascription ='abroad'").First(&node).Error
	}

	if err != nil {
		return nil, errors.New(errText)
	}
	return &node, nil
}

// 获取激活的国内国外节点数组
func (n *Node) GetEnabledNodesArr() ([]Node, error) {
	//获取节点
	var nodes []Node
	var node1 = Node{Ascription: "domestic"}
	var node2 = Node{Ascription: "abroad"}
	domesticNode, err := node1.GetEnabledNodes()
	if err != nil {
		return nil, errors.New("未正确配置活动节点")
	}
	abroadNode, err := node2.GetEnabledNodes()
	if err != nil {
		abroadNode = domesticNode
	}
	nodes = append(nodes, *domesticNode, *abroadNode)
	return nodes, nil
}

// 节点访问网站的延迟
func (n *Node) TestNodeDelayNew(url string) int64 {
	var rsp *http.Response
	var err error
	start := time.Now()
	switch runtime.GOOS {
	case "linux":
		//rsp, err = global.ClientWithDNS.Get(url)
		rsp, err = utils.ClientWithDNS("223.6.6.6", 7*time.Second).Get(url)
	default:
		//rsp, err = global.ClientWithSocks5.Get(url)
		rsp, err = utils.ClientWithSocks5("localhost", 1231, 7*time.Second).Get(url)

	}
	end := time.Since(start)
	if err != nil || rsp.StatusCode != 200 {
		return -1
	}
	defer rsp.Body.Close()
	return end.Milliseconds()
}

// 节点ip检测,t string,测试类型，domestic国内，abroad国外
func (n *Node) TestNodeIP(t string) (utils.Ip3030, error) {
	if t == "domestic" {
		return n.TestIPIPNew()
	} else {
		return n.Test3030New()
	}
}

// https://myip.ipip.net  	测试
func (n *Node) TestIPIPNew() (utils.Ip3030, error) {
	var rsp *http.Response
	var err error
	switch runtime.GOOS {
	case "linux":
		//rsp, err = global.ClientWithDNS.Get(ipTest1)
		rsp, err = utils.ClientWithDNS("223.6.6.6", 5*time.Second).Get(ipTest1)
	default:
		//rsp, err = global.ClientWithSocks5.Get(ipTest1)
		rsp, err = utils.ClientWithSocks5("localhost", 1231, 5*time.Second).Get(ipTest1)
	}
	if err != nil {
		return utils.Ip3030{}, err
	}
	//当前 IP：39.11.5.111  来自于：中国 山东 济南  联通
	r := utils.ReadDate(rsp)
	ip := strings.Split(r, "来自于：")[0]
	ip = strings.ReplaceAll(ip, "当前 IP：", "")
	ip = strings.ReplaceAll(ip, " ：", "")
	location := strings.Split(r, "来自于：")[1]
	return utils.Ip3030{ip, location}, nil
}

// http://3.0.3.0/ip   测试
func (n *Node) Test3030New() (utils.Ip3030, error) {
	var rsp *http.Response
	var err error
	switch runtime.GOOS {
	case "linux":
		//rsp, err = global.ClientWithDNS.Get(ipTest2)
		rsp, err = utils.ClientWithDNS("223.6.6.6", 5*time.Second).Get(ipTest2)
	default:
		//rsp, err = global.ClientWithSocks5.Get(ipTest2)
		rsp, err = utils.ClientWithSocks5("localhost", 1231, 5*time.Second).Get(ipTest2)
	}
	if err != nil {
		return utils.Ip3030{}, err
	}
	return utils.ReadDateToJson(rsp)
}
