package nacos_config

import "fmt"

func main() {
	client, configClient, err := NewNacosBuild().
		SetCallback(func(dataId string, context string) {

		}).
		SetClientConfigNew("AuthnAuthzServices", "ea4f46bd-61da-42cf-a9be-37ee6627a788").
		AddServer("192.168.200.100", 8848).
		AddConfigParam("mysql", "AuthnAuthzServices").
		AddRegisterService("192.168.200.203", 8081,
			"AuthnAuthzServices", "AuthnAuthzServices", "AuthnAuthzServices").
		BuildAndRun()

	fmt.Println(client, configClient, err)
}
