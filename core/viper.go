package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"mybox/core/internal"
	"mybox/global"
)

func Viper() *viper.Viper {
	config := internal.ConfigDefaultFile
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file 不存在\n"))
		} else {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.BOX_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.BOX_CONFIG); err != nil {
		fmt.Println(err)
	}
	return v
}
