package configuration

type Config struct {
	PROCESSING_WORKERS_COUNT int        `env:"PROCESSING_WORKERS_COUNT,required"`
	AMQP                     AMQPConfig `env:",prefix=AMQP_"`
}

type AMQPConfig struct {
	Host      string `env:"HOST,required"`
	Port      int    `env:"PORT,required"`
	User      string `env:"USER,required"`
	Pass      string `env:"PASS,required"`
	QueueName string `env:"QUEUE_NAME,required"`
}
