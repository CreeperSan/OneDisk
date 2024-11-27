package config

import (
	string2 "OneDisk/lib/format/formatstring"
)

// Server
// 应用配置 - 服务器配置
type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (s Server) String() string {
	return string2.String("Server{Host=%s, Port=%d}", s.Host, s.Port)
}

// Database
// 应用配置 - 数据库配置
type Database struct {
	Type     string `yaml:"type"`
	Password string `yaml:"password"`
	Path     string `yaml:"path"`
}

func (d Database) String() string {
	return string2.String("Database{Type=%s, Path=%s}", d.Type, d.Path)
}

// AppConfig
// 应用配置
type AppConfig struct {

	// 服务器信息配置
	Server Server `yaml:"server"`

	// 数据库
	Database Database `yaml:"database"`
}

func (a AppConfig) String() string {
	return string2.String("AppConfig{Server=%s, Database=%s}", a.Server, a.Database)
}
