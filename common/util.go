package common

import (
	"io/ioutil"
	"net"
	"os"
	"strings"
)

// GetLocalIPByName 通过名称获取本机Ip
func GetLocalIPByName(name string) (string, error) {
	iface, err := net.InterfaceByName("以太网")
	if err != nil {
		return "", err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	var ip string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip, nil
}

// GetLocalIP 获取本机IP
func GetLocalIP() ([]string, error) {
	ips := make([]string, 0)
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return nil, err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}

// PathExist 判断文件是否存在
func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ReadFile 根据文件地址读取文件内容
func ReadFile(path string) (string, error) {
	exist, err := PathExist(path)
	if err != nil {
		return "", err
	}

	if !exist {
		return "", nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

// WriteFile 将内容写入指定文件
func WriteFile(path, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0666)
}
