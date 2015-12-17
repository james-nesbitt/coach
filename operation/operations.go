package operation 

const (
  DEFAULT_OPERATION="<default>"
)

// Operation selector
func GetOperation() Operation {
  return Operation(&EmptyOperation{


	})
}

// Returns is a string is a valid name for an operation
func IsValidOperationName(name string) bool {
	switch name {
	case "attach","build","clean","commit","create","destroy","help","info","init","pause","pull","remove","run","scale","start","stop","tool","unpause","up":
		return true
	default:
		return false
	}
}

// Operation interface
type Operation interface {
  Run()	
}

// EmptyOperation handler for when no valid operation can be found
type EmptyOperation struct {

}
// Run method for EmptyOperation Operation interface compliance
func (operation *EmptyOperation) Run() {

}
