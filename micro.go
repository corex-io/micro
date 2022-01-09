package micro

// Service is an interface that wraps the lower level libraries
// within go-micro. Its a convenience method for building
// and initialising services.
type Service interface {
	// The service name
	Name() string
	// prepare initialises options
	Prepare(...Option)
	// Options returns the current options
	Options() Options
	// Client is used to call services
	// Client() client.Client
	// Server is for handling requests and events
	// Server() server.Server
	// Run the service
	Run() error

	Regist(Runnable)
	RegistFunc(RunFunc)
	//RegistLoop(time.Duration, Runnable)
	//RegistLoopFunc(time.Duration, RunFunc)
}

// New creates and returns a new Service based on the packages within.
func New(opts ...Option) Service {
	return newService(opts...)
}
