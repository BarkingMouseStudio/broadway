package broadway

type Config struct {
	Mailbox MailboxConfig
	Logging LoggingConfig
}

func NewConfig() Config {
	return Config{
		Mailbox: NewMailboxConfig(),
		Logging: NewLoggingConfig(),
	}
}
