package domain

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	// Agrega otros métodos según sea necesario
}
