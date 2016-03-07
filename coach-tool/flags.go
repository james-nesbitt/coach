package main

import (
	"github.com/james-nesbitt/coach/conf"
)

/**
 * Parse command flags to configure the operation
 *
 * 1: GLOBAL FLAGS : only those which we recognize below
 * 2: OPERATION [optional] : the first non-global flag, if we recognize it
 * 3. OPERATION ARGUMENTS : anything left
 *
 */
func parseGlobalFlags(flags []string) (globalFlags map[string]string, operationFlags []string, environment string) {

	globalFlags = map[string]string{} // start of with no flags

	environment = conf.COACH_CONF_ENVIRONMENTS_DEFAULT

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

		default:

			/**
			* The first flags that we don't recognize as global, fall into three cases:
			*  :{flag} : indicates an environment
			 */

			switch arg[0:1] {
			case ":": // environment
				environment = arg[1:]
			default:
				global = false
			}

		}

		// all remaining flags are local
		if !global {
			operationFlags = flags[index:]
			break
		}
	}

	// return is handles via named arguments
	return
}
