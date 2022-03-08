package data_store

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/openconfig/gnmi/proto/gnmi"
)

// type Config struct {
// 	DevicesWithMonitoring []struct {
// 		Device         interface{} `yaml:"device"`
// 		DeviceCounters []struct {
// 			Name     string `yaml:"name"`
// 			Interval int    `yaml:"interval"`
// 			Path     string `yaml:"path"`
// 		} `yaml:"device_counters"`
// 	} `yaml:"devices_with_monitoring"`
// }

func GetConfig(req *gnmi.GetRequest) []byte {
	filename, _ := filepath.Abs("./monitoring_config/example_config.yaml")
	// fmt.Println(filename)
	config_file, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Failed to read file")
	}

	// var config Config
	// fmt.Println(yaml.Unmarshal(config_file, &config))

	// fmt.Println(config)

	// return []byte(fmt.Sprintf("%v", config))

	return config_file
}
