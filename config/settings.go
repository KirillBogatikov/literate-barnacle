package config

type Settings struct {
	Port     int    `yaml:"port"`
	Postgres string `yaml:"postgres"`
	JWT      JWT    `yaml:"jwt"`
}

type JWT struct {
	PrivateKey string `yaml:"private"`
	PublicKey  string `yaml:"public"`
}
