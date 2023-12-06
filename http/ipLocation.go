package http

import (
	"encoding/json"
	"github.com/ifeisier/go-utils/text"
	"strings"
)

type IPLocation struct {
	IP        string
	Addr      string
	Longitude float64 // 经度
	Latitude  float64 // 纬度
}

func GetIPLocation(ip string) IPLocation {
	_, body, _ := GetRequest("http://whois.pconline.com.cn/ipJson.jsp?ip=" + ip + "&json=true")
	utf8, _ := text.GbkToUtf8(body)

	var bodyMap map[string]interface{}
	_ = json.Unmarshal(utf8, &bodyMap)

	_, body, _ = GetRequest("https://api.map.baidu.com/geocoder?address=" + bodyMap["city"].(string) + "&output=json")
	_ = json.Unmarshal(body, &bodyMap)

	location := IPLocation{}
	location.IP = bodyMap["ip"].(string)
	location.Addr = strings.Replace(bodyMap["addr"].(string), " ", "", -1)

	if location.Addr == "局域网" || location.Addr == "本机地址" {
		return location
	}

	lt := bodyMap["result"].(map[string]interface{})["location"].(map[string]interface{})
	location.Longitude = lt["lng"].(float64)
	location.Latitude = lt["lat"].(float64)

	return location
}
