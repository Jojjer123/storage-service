package types

type Conf struct {
	Counter []DeviceCounters `yaml:"device_counters"`
}

type DeviceCounters struct {
	Name     string `yaml:"name"`
	Interval int    `yaml:"interval"`
	Path     string `yaml:"path"`
}

type ConfigRequest struct {
	DeviceIP   string `yaml:"device_ip"`
	DeviceName string `yaml:"device_name"`
	Protocol   string `yaml:"protocol"`
	Configs    []Conf `yaml:"configs"`
}
