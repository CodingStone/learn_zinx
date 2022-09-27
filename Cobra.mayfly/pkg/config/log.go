package config

type LogFile struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Log struct {
	Level string   `yaml:"level"`
	File  *LogFile `yaml:"file"`
}
