package yamlconfig

type YmlConfigInterf interface {
	ConfigFileChangeListen()
	//容器中使用sidercar方式重新加载配置文件
	ReloadCache()
	Get(keyName string) interface{}
	GetString(keyName string) string
	GetBool(keyName string) bool
	GetStringSlice(keyNmae string) []string
}
