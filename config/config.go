package config

type Server struct {
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Zap    Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Redis  Redis  `mapstructure:"redis" json:"redis" yaml:"redis"`
	Local  Local  `mapstructure:"local" json:"local" yaml:"local"`
	JWT    JWT    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}
