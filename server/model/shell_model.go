package model

import (
	"os"
	"os/exec"
	"runtime"
	"server/global"
	"strings"
)

type Shell struct {
	Shell   string `json:"shell"`
	OutType bool   `json:"out_type"`
}

const (
	startXrayDarwin  string = "./xray -config ./config.json >/dev/null 2>&1 &"
	stopXrayDarwin          = "killall xray"
	xrayStatusDarwin        = "ps -axc | grep xray | cut -d ' ' -f1"
)

// 获取全部包
func (s *Shell) GetAllPackages() ([]string, error) {
	out, err := s.DoShell(allPackages, true)
	if err != nil {
		return nil, err
	}
	p := string(out)
	p = strings.ReplaceAll(p, "package:", "")
	pArr := strings.Split(p, "\n")
	return pArr, nil
	//fmt.Println(string(out))
}

// 包名处理
func (s *Shell) HandlePackages(packagesArr []string) string {

	var uidArr []string
	for _, v := range packagesArr {
		shStr := strings.ReplaceAll(findUid, "packageReplace", strings.TrimSpace(v))
		out, _ := s.DoShell(shStr, true)
		if len(out) > 2 {
			uidArr = append(uidArr, strings.TrimSpace(out))
		}
	}
	str := strings.Join(uidArr, " ")
	global.Logrus.Info("uid:", str)
	return str
}

// 运行shell脚本，outType bool (true:显示标准输出;false不显示标准输出)
func (s *Shell) DoShell(sh string, outType bool) (string, error) {
	os.WriteFile("temp.sh", []byte(sh), 0777)
	//defer os.Remove("temp.sh")
	var out []byte
	var err error
	var name string
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		name = "/system/bin/sh"
	default:
		name = "bash"
	}
	switch outType {
	case true:
		//fmt.Println("DoShell:", sh)
		out, err = exec.Command(name, "temp.sh").Output()
	case false:
		cmd = exec.Command(name, "temp.sh")
		err = cmd.Run()
	}
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// 查询进程状态
func (s *Shell) GetProcessStatus() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return s.DoShell(xrayStatus, true)
	default:
		return s.DoShell(xrayStatusDarwin, true)
	}
}

// 启动服务
func (s *Shell) StartService() error {
	//判断节点池工作模式
	switch global.Config.NodePoolModel {
	case "am": //开启自动切换节点
		global.NodeAutoChangeCrontab.Start()
	}
	//生成配置
	err := GenerateConfig()
	if err != nil {
		return err
	}
	//启动xray
	switch runtime.GOOS {
	case "linux":
		_, err = s.StartXrayAndroid()
		return err
	default:
		_, err = s.StartXrayDarwin()
		return err
	}
}

// 关闭服务
func (s *Shell) StopService() (string, error) {
	//判断节点池工作模式
	switch global.Config.NodePoolModel {
	case "am":
		global.NodeAutoChangeCrontab.Stop() //关闭自动切换节点
	}
	switch runtime.GOOS {
	case "linux":
		return s.StopXrayAndroid()
	default:
		return s.StopXrayDarwin()
	}
}

// 开启xray darwin
func (s *Shell) StartXrayDarwin() (string, error) {
	return s.DoShell(startXrayDarwin, false)
}

// 关闭xray darwin
func (s *Shell) StopXrayDarwin() (string, error) {
	return s.DoShell(stopXrayDarwin, false)
}

// 开启xray Android
func (s *Shell) StartXrayAndroid() (string, error) {
	//开启xray
	_, err := s.DoShell(startXray, false)
	if err != nil {
		s.DoShell(stopXray, false)
		global.Logrus.Error("StartXrayAndroid,startXray错误:", err.Error())
		return "", err
	}
	//开启rules
	config := Config{}
	config.GetConfig()
	AllowAppsStr := s.HandlePackages(config.AllowApps)
	var r string
	r = strings.ReplaceAll(startRules, "wifiProxyReplace", config.WIFIProxy)
	//r = strings.ReplaceAll(r, "IPV6NetReplace", config.IPV6Net)
	r = strings.ReplaceAll(r, "allowOutsideTcpUdpReplace", config.AllowOutsideTcpUdp)
	r = strings.ReplaceAll(r, "allowAppsUidReplace", AllowAppsStr)
	//执行规则
	_, err = s.DoShell(r, false)
	if err != nil {
		os.WriteFile("start-rules-temp.sh", []byte(r), 0777)
		s.DoShell(stopXray, false)
		s.DoShell(clearRules, false)
		global.Logrus.Error("StartXrayAndroid,startRules错误:", err.Error())
		return "", err
	}
	return "", nil
}

// 关闭xray Android
func (s *Shell) StopXrayAndroid() (string, error) {
	//关闭xray
	s.DoShell(stopXray, false)
	//清除rules
	s.DoShell(clearRules, false)
	return "", nil
}
