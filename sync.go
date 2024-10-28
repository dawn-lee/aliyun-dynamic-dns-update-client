package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// IPResponse 结构体用于解析 ifconfig.me/all.json 返回的 JSON 数据。
type IPResponse struct {
	IP string `json:"ip_addr"`
}

const (
	Endpoint = "alidns.cn-hangzhou.aliyuncs.com"
	Lang     = "en"
	Type     = "A"
	RR       = "@"
)

func CreateClient(ducConfig *Config) (_result *alidns20150109.Client, _err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(ducConfig.AccessKeyId),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		// AccessKeySecret: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
		AccessKeySecret: tea.String(ducConfig.AccessKeySecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Alidns
	config.Endpoint = tea.String(Endpoint)
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

func updateDDNS(config *Config, currentIP string) (_err error) {
	client, _err := CreateClient(config)
	if _err != nil {
		return _err
	}

	updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{
		Lang:     tea.String(Lang),
		Value:    tea.String(currentIP),
		Type:     tea.String(Type),
		RR:       tea.String(RR),
		RecordId: tea.String("1"),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err = client.UpdateDomainRecordWithOptions(updateDomainRecordRequest, runtime)
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 此处仅做打印展示，请谨慎对待异常处理，在工程项目中切勿直接忽略异常。
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return _err
}

// getPublicIP 获取当前的公网IP地址。
func getPublicIP() (string, error) {
	resp, err := http.Get("http://ifconfig.me/all.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ipResp IPResponse
	err = json.Unmarshal(body, &ipResp)
	return ipResp.IP, err
}
