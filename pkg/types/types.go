package types

type Conf struct {
	Counter []DeviceCounters `yaml:"device_counters"`
}

type DeviceCounters struct {
	Name     string `yaml:"name"`
	Interval int    `yaml:"interval"`
	Path     string `yaml:"path"`
}

type ConfigObject struct {
	DeviceIP   string `yaml:"device_ip"`
	DeviceName string `yaml:"device_name"`
	Protocol   string `yaml:"protocol"`
	Configs    []Conf `yaml:"configs"`
}

// type DevicesWithMonitoring struct {
// 	DeviceIP   string `yaml:"device_ip"`
// 	DeviceName string `yaml:"device_name"`
// 	Protocol   string `yaml:"protocol"`
// 	Configs    []Conf
// 	} `yaml:"configs"`
// }

type Config struct {
	Devices []ConfigObject `yaml:"devices_with_monitoring"`
}

// Will most likely be used to store the config in storage
type Elem struct {
	Namespace string
	Name      string
	Value     string
	Elem      *Elem
}

type Module struct {
	Elements []Elem
}

type SchemaTree struct {
	Modules []Module
}