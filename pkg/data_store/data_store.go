package data_store

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/openconfig/gnmi/proto/gnmi"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DevicesWithMonitoring []struct {
		DeviceIP   string `yaml:"device_ip"`
		DeviceName string `yaml:"device_name"`
		Protocol   string `yaml:"protocol"`
		Configs    []struct {
			DeviceCounters []struct {
				Name     string `yaml:"name"`
				Interval int    `yaml:"interval"`
				Path     string `yaml:"path"`
			} `yaml:"device_counters"`
		} `yaml:"configs"`
	} `yaml:"devices_with_monitoring"`
}

func GetConfig(req *gnmi.GetRequest) []byte {
	filename, _ := filepath.Abs("./monitoring_config/example_config.yaml")
	// fmt.Println(filename)
	config_file, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Failed to read file")
	}

	var config Config
	yaml.Unmarshal(config_file, &config)

	// fmt.Println(config)

	index := -1
	for i, device := range config.DevicesWithMonitoring {
		if device.DeviceIP == req.Path[0].Target {
			index = i
			break
		}
	}

	var conf []byte

	if index != -1 {
		conf, err = yaml.Marshal(config.DevicesWithMonitoring[index])
	} else {
		fmt.Println("Could not find the device in config!")
	}

	if err != nil {
		fmt.Println("Failed to convert config to byte slice")
	}

	return conf

	// test := []byte(fmt.Sprintf("%v", config))

	// var conf Config
	// yaml.Unmarshal(test, &conf)
	// fmt.Println(conf.DevicesWithMonitoring[0])

	// fmt.Println(config)

	// return []byte(fmt.Sprintf("%v", config))

	// return nil
	// return config_file
}
