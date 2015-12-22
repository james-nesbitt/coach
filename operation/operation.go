package operation

const (
	DEFAULT_OPERATION = "<default operation>"
)

// Configure a set of operations
type OperationsSettings struct {
}

// A set of Operations
type Operations struct {
	OperationsSettings
}

// Validate a string as an operation name
func IsValidOperationName(arg string) bool {
	return true
}

// Operation that can act on A target list
type Operation interface {
}
