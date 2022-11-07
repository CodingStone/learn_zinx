package config

import "fmt"

type Tls struct {
	Enable   bool   `yaml:"enable"`   //是否开启
	KeyFile  string `yaml:"key-file"` //私钥位置
	CertFile string `yaml:"cer-file"` //证书位置
}

type Static struct {
	RelativePath string `yaml:"relative-path"`
	Root         string `yaml:"root"`
}

type StaticFile struct {
	RelativePath string `yaml:"relative-path"`
	Filepath     string `yaml:"filepath"`
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

// 获取终端回访记录存放基础路径, 如果配置文件未配置，则默认为./rec
func (s *Server) GetMachineRecPath() string {
	path := s.MachineRecPath
	if path == "" {
		return "./rec"
	}
	return path
}

func (s *Server) GetPort() string {
	return fmt.Sprintf(":%d", s.Port)
}
