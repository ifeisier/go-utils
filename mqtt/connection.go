package mqtt

import (
    "context"
    "github.com/eclipse/paho.golang/autopaho"
    "github.com/eclipse/paho.golang/paho"
    "go-linux-iot/logx"
    "go-linux-iot/shared"
    "net/url"
)

// SetClientID 设置客户端 ID
//
// clientID: 客户端 ID
func (c *Client) SetClientID(clientID string) {
    c.clientID = clientID
}

// SetUsernamePassword 设置用户名和密码
//
// un: 用户名
//
// pwd: 密码
func (c *Client) SetUsernamePassword(un, pwd string) {
    c.username = un
    c.password = []byte(pwd)
}

// Connection 没有连接会创建新的连接，否则就会返回已经创建的连接
func (c *Client) Connection() (*autopaho.ConnectionManager, error) {
    c.once.Do(func() {
        urls := make([]*url.URL, 0, len(c.serverURL))
        for _, v := range c.serverURL {
            su, _ := url.Parse(v)
            urls = append(urls, su)
        }

        cliCfg := autopaho.ClientConfig{
            BrokerUrls:                    urls,
            KeepAlive:                     30,
            CleanStartOnInitialConnection: true,
            SessionExpiryInterval:         0,
            OnConnectionUp:                c.connectionUp,
            OnConnectError:                c.connectError,
            Debug:                         logx.PahoMQTTLog{Prefix: "MQTTDebug", IsErrorLog: false},
            Errors:                        logx.PahoMQTTLog{Prefix: "MQTTErrors", IsErrorLog: true},
            PahoDebug:                     logx.PahoMQTTLog{Prefix: "MQTTPahoDebug", IsErrorLog: false},
            PahoErrors:                    logx.PahoMQTTLog{Prefix: "MQTTPahoErrors", IsErrorLog: true},
            ClientConfig: paho.ClientConfig{
                ClientID:           c.clientID,
                OnClientError:      c.clientError,
                OnServerDisconnect: c.serverDisconnect,
            },
        }

        cliCfg.ConnectUsername = c.username
        cliCfg.ConnectPassword = c.password

        cm, err := autopaho.NewConnection(context.Background(), cliCfg)
        if err != nil {
            c.connectionManager = cm
            c.error = err
            return
        }

        err = cm.AwaitConnection(context.Background())
        if err != nil {
            c.connectionManager = nil
            c.error = err
            return
        }
        c.connectionManager = cm
    })

    return c.connectionManager, c.error
}

// Publish 发布消息
func (c *Client) Publish(ctx context.Context, p *paho.Publish) (*paho.PublishResponse, error) {
    return c.connectionManager.Publish(ctx, p)
}

// RegisterHandler 注册消息处理函数
func (c *Client) RegisterHandler(topic string, handler paho.MessageHandler) {
    c.clientConfig.Router.RegisterHandler(topic, handler)
}

// connectionUp 与 MQTT 服务器连接成功（包括重新连接）时，会回调这个函数。
func (c *Client) connectionUp(cm *autopaho.ConnectionManager, pc *paho.Connack) {
    shared.LOGGER.Info("连接到 MQTT 服务端")

    if _, err := cm.Subscribe(context.Background(), c.subscribe); err != nil {
        m := make(map[string]any)
        m["err"] = err.Error()
        shared.LOGGER.Error("订阅失败:", logx.MsgInfo(m))
    }
}

// connectError 连接失败时，会回调这个函数。
func (c *Client) connectError(err error) {
    m := make(map[string]any)
    m["err"] = err.Error()
    shared.LOGGER.Error("尝试连接到 MQTT 客户端出错:", logx.MsgInfo(m))
}

func (c *Client) clientError(err error) {
    m := make(map[string]any)
    m["err"] = err.Error()
    shared.LOGGER.Error("网络原因导致连接断开:", logx.MsgInfo(m))
}

// serverDisconnect 收到 MQTT 服务器 packets.DISCONNECT 时，会回调这个函数。
func (c *Client) serverDisconnect(d *paho.Disconnect) {
    m := make(map[string]any)
    m["reason"] = d.Properties.ReasonString
    m["code"] = d.ReasonCode
    shared.LOGGER.Error("收到服务端的断开连接命令:", logx.MsgInfo(m))
}
