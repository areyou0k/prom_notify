package yamlconfig

import (
	"log"
	"promnotify/container"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var lastChangeTime time.Time
var fileName = "alarmgroup"
var ConfigKeyPrefix = "Config_"
var ConfigYml YmlConfigInterf

func init() {
	lastChangeTime = time.Now()
}

// 创建一个yaml配置文件工厂
// 参数设置为可变参数的文件名，这样参数就可以不需要传递，如果传递了多个，我们只取第一个参数作为配置文件名
func CreateYamlFactory() {

	yamlConfig := viper.New()
	// yamlConfig.AddConfigPath("./")
	yamlConfig.AddConfigPath("/etc/app")
	yamlConfig.SetConfigName(fileName)
	yamlConfig.SetConfigType("yaml")

	if err := yamlConfig.ReadInConfig(); err != nil {
		log.Fatal("load yaml config err: ", err)
	}

	// 全局配置指针
	ConfigYml = &ymlConfig{
		yamlConfig,
	}
	// return &ymlConfig{
	// 	yamlConfig,
	// }
}

type ymlConfig struct {
	viper *viper.Viper
}

//监听文件变化
func (y *ymlConfig) ConfigFileChangeListen() {
	y.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Since(lastChangeTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				log.Printf("OnConfigChange: %s Op:%s\n", changeEvent.Name, changeEvent.Op)
				y.clearCache()
				lastChangeTime = time.Now()
			}
		}
	})
	y.viper.WatchConfig()
}

// 判断相关键是否已经缓存
func (y *ymlConfig) keyIsCache(keyName string) bool {
	if _, exists := container.CreateContainersFactory().KeyIsExists(ConfigKeyPrefix + keyName); exists {
		return true
	} else {
		return false
	}
}

// 对键值进行缓存
func (y *ymlConfig) cache(keyName string, value interface{}) bool {
	return container.CreateContainersFactory().Set(ConfigKeyPrefix+keyName, value)
}

// 通过键获取缓存的值
func (y *ymlConfig) getValueFromCache(keyName string) interface{} {
	return container.CreateContainersFactory().Get(ConfigKeyPrefix + keyName)
}

// 清空已经窜换的配置项信息
func (y *ymlConfig) clearCache() {
	container.CreateContainersFactory().FuzzyDelete(ConfigKeyPrefix)
}

// 重新加载配置文件
func (y *ymlConfig) ReloadCache() {
	y.clearCache()
	CreateYamlFactory()
}

// Get 一个原始值
func (y *ymlConfig) Get(keyName string) interface{} {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName)
	} else {
		value := y.viper.Get(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetString
func (y *ymlConfig) GetString(keyName string) string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(string)
	} else {
		value := y.viper.GetString(keyName)
		y.cache(keyName, value)
		return value
	}

}

// GetBool
func (y *ymlConfig) GetBool(keyName string) bool {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(bool)
	} else {
		value := y.viper.GetBool(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetStringSlice
func (y *ymlConfig) GetStringSlice(keyName string) []string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).([]string)
	} else {
		value := y.viper.GetStringSlice(keyName)
		y.cache(keyName, value)
		return value
	}
}
