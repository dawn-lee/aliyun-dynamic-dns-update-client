package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Permit  Permit   `json:"Permit"`  // 阿里云通行证
	Domains []Domain `json:"Domains"` // 域名列表
}

type Permit struct {
	AccessKeyId     string `json:"AccessKeyId"`     // 阿里云AccessKeyId
	AccessKeySecret string `json:"AccessKeySecret"` // 阿里云AccessKeySecret
}

type Domain struct {
	Domain   string `json:"Domain"`   // 域名
	RecordId string `json:"RecordId"` // 域名ID
	Lang     string `json:"Lang"`     // 语言
	Type     string `json:"Type"`     // 解析类型
	RR       string `json:"RR"`       // 主机记录
}

const (
	lastIpFileName = "currentIp.txt"
	configFileName = "config.json"
)

// setLastIP 将新的IP地址写入文件。
func setLastIP(ip string) error {
	return ioutil.WriteFile(lastIpFileName, []byte(ip), 0644)
}

// 保存app配置
func saveConfig(config *Config) error {
	// 将 config 结构体序列化为 JSON 字符串
	configBytes, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		return err
	}
	ioutil.WriteFile(configFileName, []byte(configBytes), 0644)
	return nil
}

// 保存阿里云通行证 ak,sk
func savePermit(permit *Permit) error {
	// 将 permit 结构体序列化为 JSON 字符串
	// configBytes, err := json.MarshalIndent(permit, "", "	")
	// if err != nil {
	// 	return err
	// }
	return nil
}

// 读取app配置
func readConfig() (Config, error) {
	body, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(body, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

// 读取阿里云通行证 ak,sk
func readPermit() (Permit, error) {
	config, _err := readConfig()
	if _err != nil {
		return Permit{}, _err
	}
	return config.Permit, nil
}

// 获取所有domain
func readDomains() ([]Domain, error) {
	config, _err := readConfig()
	if _err != nil {
		return []Domain{}, _err
	}
	return config.Domains, nil
}
