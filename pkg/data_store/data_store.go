package data_store

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"encoding/json"

	"github.com/openconfig/gnmi/proto/gnmi"
	"gopkg.in/yaml.v2"

	types "github.com/onosproject/storage-service/pkg/types"
)

func GetConfig(req *gnmi.GetRequest) []byte {
	filename, _ := filepath.Abs("./monitoring_config/example_config.yaml")
	// fmt.Println(filename)
	config_file, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Failed to read file")
	}

	var config types.Config
	yaml.Unmarshal(config_file, &config)

	// fmt.Println(config)

	index := -1
	for i, device := range config.Devices {
		if device.DeviceIP == req.Path[0].Target {
			index = i
			break
		}
	}

	var conf []byte

	if index != -1 {
		conf, err = yaml.Marshal(config.Devices[index])
	} else {
		fmt.Println("Could not find the device in config!")
	}

	if err != nil {
		fmt.Println("Failed to convert config to byte slice")
	}

	return conf
}

func GetAdapter(req *gnmi.GetRequest) []byte {
	switch req.Path[0].Target {
	case "NETCONF":
		adapter := types.Adapter{
			Protocol: "NETCONF",
			Address:  "gnmi-netconf-adapter:11161",
		}
		fmt.Println(adapter)
		objSlice, err := json.Marshal(adapter)
		if err != nil {
			fmt.Println("Failed to marshal adapter!")
		}
		return objSlice
	default:
		fmt.Println("Did not recognize protocol when getting adapter!")
	}

	return nil
}

func GetFullConfig() types.Config {
	filename, _ := filepath.Abs("./monitoring_config/example_config.yaml")
	config_file, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Failed to read file")
	}

	var config types.Config
	yaml.Unmarshal(config_file, &config)

	return config
}

func UpdateConfig(conf types.Config) error {
	// filename, _ := filepath.Abs("./monitoring_config/example_config.yaml")
	// fmt.Println(filename)
	// config_file, err := ioutil.ReadFile(filename)

	// if err != nil {
	// 	fmt.Println("Failed to read file")
	// }

	// var config Config
	// yaml.Unmarshal(config_file, &config)

	// // fmt.Println(config)

	// index := -1
	// for i, device := range config.DevicesWithMonitoring {
	// 	if device.DeviceIP == req.Path[0].Target {
	// 		index = i
	// 		break
	// 	}
	// }

	return nil
}
