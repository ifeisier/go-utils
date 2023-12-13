package http

import (
	"encoding/json"
	"fmt"
	"github.com/ifeisier/go-utils/text"
	"strings"
)

type IPLocation struct {
	IP        string
	Addr      string
	Longitude float64 // 经度
	Latitude  float64 // 纬度
}

func GetIPLocation(ip string) (l *IPLocation, e error) {
	defer func() {
		if err := recover(); err != nil {
			l = nil
			e = fmt.Errorf("未知错误")
		}
	}()

	_, body, err := GetRequest("http://whois.pconline.com.cn/ipJson.jsp?ip=" + ip + "&json=true")
	if err != nil {
		return nil, err
	}
	utf8, err := text.GbkToUtf8(body)
	if err != nil {
		return nil, err
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(utf8, &bodyMap)
	if err != nil {
		return nil, err
	}

	_, body, err = GetRequest("https://api.map.baidu.com/geocoder?address=" + bodyMap["city"].(string) + "&output=json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	location := &IPLocation{}
	location.IP = bodyMap["ip"].(string)
	location.Addr = strings.Replace(bodyMap["addr"].(string), " ", "", -1)

	if location.Addr == "局域网" || location.Addr == "本机地址" {
		return location, nil
	}

	lt := bodyMap["result"].(map[string]interface{})["location"].(map[string]interface{})
	location.Longitude = lt["lng"].(float64)
	location.Latitude = lt["lat"].(float64)

	return location, nil
}
