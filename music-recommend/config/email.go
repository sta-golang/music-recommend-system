package config

type EmailConfig struct {
	Host  string `json:"host" yaml:"host"`
	Port  int    `json:"port" yaml:"port"`
	Email string `json:"email" yaml:"email"`
	Pwd   string `json:"pwd" yaml:"pwd"`
}
