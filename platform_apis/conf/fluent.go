package conf

type Fluent struct {
	Host string `mapstructure:"fluent_host"`
	Port int    `mapstructure:"fluent_port"`
}
