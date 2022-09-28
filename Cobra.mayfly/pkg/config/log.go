package config

import "path"

type LogFile struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Log struct {
	Level string   `yaml:"level"`
	File  *LogFile `yaml:"file"`
}

// 获取完整路径文件名[这种初始化后，然后直接带方法，go设计真的很好]
func (l *LogFile) GetFilename() string {
	var filepath, filename string
	if fp := l.Path; fp == "" {
		filepath = "./"
	} else {
		filepath = fp
	}
	if fn := l.Name; fn == "" {
		filename = "default.log"
	} else {
		filename = fn
	}

	return path.Join(filepath, filename)
}
