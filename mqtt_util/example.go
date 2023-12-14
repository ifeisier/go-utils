package mqtt_util

import (
	"github.com/eclipse/paho.golang/autopaho"
	"math/rand"
	"strconv"
	"time"
)

func initMQTT() (*autopaho.ConnectionManager, error) {
	clientID := "api-"
	int31n := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(16)
	for _, v := range int31n {
		clientID += strconv.Itoa(v)
	}

	manager, err := NewClientBuild().AddServerURL("mqtt://192.168.200.100:1883").
		AddSubscribe("$share/IOTDataServices/+/+/thingsModel/up").SetClientID(clientID).
		SetUsernamePassword("fE^8*GFs@kbr6P5!pIOUMjyIwj&G#iax", "wQYSH9RkT4MAMR4CZ88lQ%nZ7H6r7bx2").
		BuildAndConnection()
	return manager, err
}
