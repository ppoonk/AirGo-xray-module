package model

import (
	"errors"
	"server/global"
	"server/utils"
	"strings"
	"time"
)

type Subscription struct {
	//gorm.Model
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Url   string `json:"url"`
	Alias string `json:"alias"`
	Nodes []Node `json:"nodes" gorm:"foreignKey:SubscriptionID"`
}

// 解析订阅
func (ss *Subscription) ParseSub() error {
	//去掉前后空格
	ss.Url = strings.Trim(ss.Url, " \n")
	var linkArr []string
	if strings.HasPrefix(ss.Url, "http") { //sub
		linkArr1, err := ss.ParseUrl()
		if err != nil {
			return err
		}
		linkArr = linkArr1
	} else if strings.HasPrefix(ss.Url, "v") || strings.HasPrefix(ss.Url, "t") { //node
		//系统订阅
		ss.ID = 1
		linkArr = ss.ParseBase64(ss.Url)
	} else {
		//fmt.Println("订阅格式:", ss.Url)
		return errors.New("未知订阅格式")
	}
	for _, v := range linkArr {
		data := ParseLink(v)
		if data == nil {
			continue
		}
		ss.Nodes = append(ss.Nodes, *data)
	}
	return ss.AddSubToDb() //has many自动创建node
}

// 更新订阅
func (ss *Subscription) UpdateSub() error {
	//find
	err := global.DB.First(&ss, ss.ID).Error
	if err != nil {
		return err
	}
	//
	linkArr, err := ss.ParseUrl()
	if err != nil {
		return err
	}
	for _, v := range linkArr {
		data := ParseLink(v) //解析每一条订阅
		ss.Nodes = append(ss.Nodes, *data)
	}
	ss.DeleteSubNode()
	return ss.ReplaceSubNode()
}

// 解析订阅url
func (ss *Subscription) ParseUrl() ([]string, error) {
	var linkArr []string
	//rsp, err := global.ClientWithDNS.Get(ss.Url)
	rsp, err := utils.ClientWithDNS("223.6.6.6", 5*time.Second).Get(ss.Url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close() ///
	subLink := utils.ReadDate(rsp)
	//rsp.Body.Close()
	if len(subLink) == 0 {
		return nil, errors.New("节点列表为空")
	}
	linkArr = ss.ParseBase64(subLink)
	return linkArr, nil
}

// 解析订阅文本
func (ss *Subscription) ParseBase64(subtext string) []string {
	var data string
	//如果传进来vmess vless trojan开头的多个节点，则跳过base64解码
	if strings.HasPrefix(ss.Url, "v") || strings.HasPrefix(ss.Url, "t") {
		data = subtext
	} else {
		data = utils.SubBase64Decode(subtext)
	}
	s := strings.ReplaceAll(data, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	list := strings.Split(strings.TrimRight(s, "\n"), "\n")
	return list
}

// 添加订阅
func (ss *Subscription) AddSubToDb() error {
	//系统订阅
	if ss.ID == 1 {
		return global.DB.Debug().Model(&Subscription{ID: ss.ID}).Association("Nodes").Append(&ss.Nodes)
	}
	return global.DB.Debug().Create(&ss).Error
}

// 删除订阅
func (ss *Subscription) DeleteSub() error {
	err := ss.DeleteSubNode()
	if err != nil {
		return err
	}
	return global.DB.Debug().Delete(&ss).Error
}

// 替换该订阅关联的节点
func (ss *Subscription) ReplaceSubNode() error {
	return global.DB.Model(&ss).Association("Nodes").Replace(ss.Nodes)
}

// 删除该订阅关联的节点
func (ss *Subscription) DeleteSubNode() error {
	return global.DB.Model(&Node{}).Where("subscription_id = ?", ss.ID).Delete(&Node{}).Error
}

// 获取订阅关联的节点列表
func (ss *Subscription) GetSubNodeList() error {
	return global.DB.Where("subscription_id = ?", ss.ID).Find(&ss.Nodes).Error
}

// 获取订阅列表
func (ss *Subscription) GetSubList() (*[]Subscription, error) {
	var subList []Subscription
	err := global.DB.Find(&subList).Error
	return &subList, err
}
