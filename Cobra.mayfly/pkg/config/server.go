package config

type Tls struct {
	Enable  bool   `yaml:"enable"`   //是否开启
	KeyFile string `yaml:"key_file"` //私钥位置
	CerFile string `yaml:"cer_file"` //证书位置
}

type Static struct {
	RelativePath string `yaml:"relative_path"`
	Root         string `yaml:"root"`
}

type StaticFile struct {
	RelativePath string `yaml:"relative_path"`
	Root         string `yaml:"root"`
}

type Server struct {
	Port           int            `yaml:"port"`
	Model          string         `yaml:"model"`
	Cors           bool           `yaml:"cors"`
	Tls            *Tls           `yaml:"tls"`
	Static         *[]*Static     `yaml:"static"`
	StaticFile     *[]*StaticFile `yaml:"static_file"`
	MachineRecPath string         `yaml:"machine_rec_path"` //机器终端操作回放文件存储路径
}
