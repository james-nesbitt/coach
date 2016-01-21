package main

import (
	"github.com/james-nesbitt/coach/operation"
)

/**
 * Parse command flags to configure the operation
 *
 * 1: GLOBAL FLAGS : only those which we recognize below
 * 2: OPERATION [optional] : the first non-global flag, if we recognize it
 * 3. OPERATION ARGUMENTS : anything left
 *
 */
func parseGlobalFlags(flags []string) (operationName string, targetIdentifiers []string, globalFlags map[string]string, operationFlags []string) {
	operationName = operation.DEFAULT_OPERATION // default operation, to be interpreted later, if not set in this function

	globalFlags = map[string]string{} // start of with no flags
	targetIdentifiers = []string{}    //  ||

	global := true // start of assuming everything is a global arg
	for index := 1; index < len(flags); index++ {
		arg := flags[index]

		switch arg {
		case "-v":
			fallthrough
		case "--info":
			globalFlags["verbosity"] = "info"
		case "-vv":
			fallthrough
		case "--verbose":
			globalFlags["verbosity"] = "verbose"
		case "-vvv":
			fallthrough
		case "--debug":
			globalFlags["verbosity"] = "debug"
		case "-vvvv":
			fallthrough
		case "--staaap":
			globalFlags["verbosity"] = "staaap"

		case "--all": // this is default anyway
			targetIdentifiers = append(targetIdentifiers, "$all")

		default:

			/**
			* The first flags that we don't recognize as global, fall into three cases:
			*  @{flag} : indicates a node target, can be repeated
			*  %{flag} : indicates a node type target, can be repeated
			*  -{flag} : indicates the end of global flag targeting, and starts the collection of operationFlags
			*  {flag} : (first only) indicates which operation (default is info)
			 */

			switch arg[0:1] {
			case "@": // target
				fallthrough
			case "%": // type
				targetIdentifiers = append(targetIdentifiers, arg)

			// this means that local flags have started being processed, as all global flags are particular
			case "-": // local flag
				global = false

			default: // operation

				// if we recognize the subsequent argument as an operation, then set it,
				// otherwise we assume a default operation, and that op args have started
				// @TODO there has got to be a better way of doing this.
				if operation.IsValidOperationName(arg) {
					operationName = arg
					index++
				}

				global = false

			}

		}

		// all remaining flags are local
		if !global {
			operationFlags = flags[index:]
			break
		}
	}

	// if not targets were specifiec, then use all the targets
	if len(targetIdentifiers) == 0 {
		targetIdentifiers = []string{"$all"}
	}

	// return is handles via named arguments
	return
}
