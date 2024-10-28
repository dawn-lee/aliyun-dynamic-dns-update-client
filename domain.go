package main

import (
	"fmt"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
)

// 立即更新
func refresh() error {
	// 获取阿里云凭据
	permit, _ := readPermit()
	// 创建client
	client, _err := CreateClient(permit)
	if _err != nil {
		return _err
	}
	// 获取公共IP
	currentIp, _err := getPublicIP()
	if _err != nil {
		return _err
	}
	// 获取配置的所有domain
	domains, _err := readDomains()
	if _err != nil {
		return _err
	}
	for _, domain := range domains {
		result := update(domain)
		fmt.Println("是否更新：", result)
		if result {
			updateDDNS(client, domain, currentIp)
		}
	}
	return nil
}

// 是否更新
func update(domain Domain) (_result bool) {
	permit, _ := readPermit()
	// 创建client
	client, _err := CreateClient(permit)
	if _err != nil {
		return false
	}
	// 查询domain下的record
	result, _err := queryDomainRecords(domain.Domain, client)
	if _err != nil {
		return false
	}
	recordMap := make(map[string]*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord)
	for _, record := range result.Body.DomainRecords.Record {
		recordMap[*record.RecordId] = record
	}
	record := recordMap[domain.RecordId]
	currentIp, _err := getPublicIP()
	if domain.Type == *record.Type && domain.RR == *record.RR && currentIp == *record.Value {
		return false
	}
	return true
}
