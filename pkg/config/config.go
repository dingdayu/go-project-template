package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// changeEventHandle 配置变更处理器
var (
	changeEventHandle []func(e fsnotify.Event)
	eventLock         sync.Mutex
	once              sync.Once
)

// CfgFile 配置文件路径,允许在初始化前,由外部包赋值
var CfgFile string

func Init() {
	once.Do(func() {
		// 设置配置文件目录和文件名
		viper.SetConfigName("config")                 // name of config file (without extension)
		viper.SetConfigType("yaml")                   // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath("/etc/singbox-adapter/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.singbox-adapter") // call multiple times to add many search paths
		viper.AddConfigPath(".")                      // optionally look for config in the working directory

		if len(CfgFile) > 0 {
			viper.SetConfigFile(CfgFile)
		}

		// read in environment variables that match
		viper.SetEnvPrefix("GO")
		viper.AutomaticEnv()

		viper.BindEnv("app.service_name", "OTEL_SERVICE_NAME")
		viper.BindEnv("app.port", "HTTP_PORT")
		viper.BindEnv("app.environment", "ENVIRONMENT")
		viper.BindEnv("jwt.secret", "JWT_SECRET")
		viper.BindEnv("db", "DB")

		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		if err := viper.ReadInConfig(); err == nil {
			fmt.Printf("\033[1;30;42m[info]\033[0m using config file %s\n", viper.ConfigFileUsed())
		} else {
			fmt.Printf("\033[1;30;41m[error]\033[0m using config file error %s\n", err.Error())
			os.Exit(1)
		}

		// // 监听配置文件变更
		// viper.WatchConfig()
		// // 调用 config 变更注册
		// viper.OnConfigChange(onConfigChange)

		fmt.Printf("\033[1;30;42m[info]\033[0m config init %s\n", viper.ConfigFileUsed())
	})
}

// RegisterChangeEvent 注册配置变更事件
func RegisterChangeEvent(f func(e fsnotify.Event)) {
	eventLock.Lock()
	defer eventLock.Unlock()

	changeEventHandle = append(changeEventHandle, f)
}

// onConfigChange 循环执行事件调用
//
//nolint:unused
func onConfigChange(e fsnotify.Event) {
	for _, f := range changeEventHandle {
		f(e)
	}
}
