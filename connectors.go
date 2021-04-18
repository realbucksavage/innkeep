package innkeep

// Connector represents a "transport" receiver that can receive application infomation
type Connector interface {

	// A friendly name
	Name() string

	// Initialize and start the connector
	Start() error

	// Perform a graceful shutdown
	Stop() error
}
