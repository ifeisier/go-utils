package nacos_util

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type Callback func(string, string)

type NacosBuild struct {
	ServerConfigs         []constant.ServerConfig
	ClientConfig          *constant.ClientConfig
	RegisterInstanceParam []vo.RegisterInstanceParam
	ConfigParam           []vo.ConfigParam
	callback              Callback
}

func NewNacosBuild() *NacosBuild {
	return &NacosBuild{callback: func(string, string) {}}
}

// BuildAndRun 构建并运行
func (build *NacosBuild) BuildAndRun() (
	iNamingClient naming_client.INamingClient, iConfigClient config_client.IConfigClient, err error) {

	if build.ServerConfigs == nil || build.ClientConfig == nil {
		err = fmt.Errorf("没有指定 Nacos 服务器或客户端配置")
		return
	}

	if len(build.RegisterInstanceParam) != 0 {
		iNamingClient, err = build.newNamingClient()
		if err != nil {
			return
		}
	}

	if len(build.ConfigParam) != 0 {
		iConfigClient, err = build.newConfigClient()
		if err != nil {
			return
		}
	}

	return iNamingClient, iConfigClient, nil
}

func (build *NacosBuild) newConfigClient() (iClient config_client.IConfigClient, err error) {
	iClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  build.ClientConfig,
			ServerConfigs: build.ServerConfigs,
		},
	)
	if err != nil {
		return
	}

	for _, configParam := range build.ConfigParam {
		context, err := iClient.GetConfig(configParam)
		if err != nil {
			return nil, err
		}
		build.callback(configParam.DataId, context)

		// 创建配置文件监听
		configParam.OnChange = func(namespace, group, dataId, data string) {
			build.callback(dataId, data)
		}
		err = iClient.ListenConfig(configParam)
		if err != nil {
			return nil, err
		}
	}

	return iClient, nil
}

func (build *NacosBuild) newNamingClient() (iClient naming_client.INamingClient, err error) {
	iClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  build.ClientConfig,
			ServerConfigs: build.ServerConfigs,
		},
	)
	if err != nil {
		return
	}

	for _, registerInstanceParam := range build.RegisterInstanceParam {
		_, err = iClient.RegisterInstance(registerInstanceParam)
	}

	return
}

// SetCallback 设置回调函数
func (build *NacosBuild) SetCallback(cb Callback) *NacosBuild {
	build.callback = cb
	return build
}

// AddServerConfig 添加 Nacos 服务器
func (build *NacosBuild) AddServerConfig(serverConfig constant.ServerConfig) *NacosBuild {
	build.ServerConfigs = append(build.ServerConfigs, serverConfig)
	return build
}

// AddServer 添加 Nacos 服务器
func (build *NacosBuild) AddServer(ip string, port uint64) *NacosBuild {
	build.ServerConfigs = append(build.ServerConfigs, constant.ServerConfig{IpAddr: ip, Port: port})
	return build
}

// SetClientConfig 设置客户端配置
func (build *NacosBuild) SetClientConfig(clientConfig *constant.ClientConfig) *NacosBuild {
	build.ClientConfig = clientConfig
	return build
}

// SetClientConfigNew 设置客户端配置
func (build *NacosBuild) SetClientConfigNew(appName, namespaceId string) *NacosBuild {
	build.ClientConfig = NewClientConfig(appName, namespaceId)
	return build
}

// AddConfigParam 添加配置中心
func (build *NacosBuild) AddConfigParam(dataId, group string) *NacosBuild {
	build.ConfigParam = append(build.ConfigParam, vo.ConfigParam{DataId: dataId, Group: group})
	return build
}

// AddConfig 添加配置中心
func (build *NacosBuild) AddConfig(configParam vo.ConfigParam) *NacosBuild {
	build.ConfigParam = append(build.ConfigParam, configParam)
	return build
}

// AddRegisterService 注册服务
func (build *NacosBuild) AddRegisterService(localIP string, port uint64, serviceName, clusterName, groupName string) *NacosBuild {
	build.RegisterInstanceParam = append(build.RegisterInstanceParam,
		NewRegisterInstanceParam(localIP, port, serviceName, clusterName, groupName))
	return build
}

// AddRegisterInstanceService 注册服务
func (build *NacosBuild) AddRegisterInstanceService(registerInstanceParam vo.RegisterInstanceParam) *NacosBuild {
	build.RegisterInstanceParam = append(build.RegisterInstanceParam, registerInstanceParam)
	return build
}

// NewRegisterInstanceParam 创建服务注册参数
func NewRegisterInstanceParam(localIP string, port uint64, serviceName, clusterName, groupName string) vo.RegisterInstanceParam {
	return vo.RegisterInstanceParam{
		Ip:          localIP,
		Port:        port,
		ServiceName: serviceName,
		ClusterName: clusterName,
		GroupName:   groupName,
		Metadata:    map[string]string{},
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}
}

// NewClientConfig 创建客户端配置
func NewClientConfig(appName, namespaceId string) *constant.ClientConfig {
	return &constant.ClientConfig{
		AppName:              appName,
		NamespaceId:          namespaceId,
		CacheDir:             "./caches",
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: true,
		LogDir:               "./logs",
		LogLevel:             "debug",
		LogRollingConfig:     &constant.ClientLogRollingConfig{MaxSize: 10, MaxAge: 10, MaxBackups: 10, LocalTime: true},
		AppendToStdout:       true,
	}
}
