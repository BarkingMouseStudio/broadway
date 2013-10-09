package broadway

type Config struct {
	mailbox MailboxConfig
	logging LoggingConfig
}

func NewConfig() Config {
	return Config{
		mailbox: MailboxConfig{},
		logging: NewLoggingConfig(),
	}
}
