package conf

type Redis struct {
	RedisHost string `viper:"string" mapstructure:"redis_host"`
	RedisPort string `viper:"string" mapstructure:"redis_port"`
}
