package conf

/**
 * @Note that it is advisable to not test if a path exists, as 
 * that can cause race conditions, and can produce an invalid test
 * as the path could be created between the test, and the use of
 * the path.
 */

import (
 	"os"
	"path"
)

func MakePaths() Paths {
	return Paths{
		Paths: map[string]string{},
		ConfPathKeys: []string{},
	}
}

type Paths struct {
	Paths map[string]string // Paths is a string keyed map of important paths in the project
	ConfPathKeys []string // ConfPathKeys is an ordered list of Path keys, in which project configurations could be kept
}

// If a path key has been set, return it
func (paths *Paths) Path(key string) (path string, ok bool) {
	path, ok = paths.Paths[key]
	return
}
// Add a new Path to the Paths
func (paths *Paths) SetPath(key string, keyPath string, overwrite bool) bool {
	if _, ok := paths.Paths[key]; overwrite || !ok {
		paths.Paths[key] = keyPath
	}
	return false
}
// Set a path key as a ConfPath
func (paths *Paths) setConfPath(key string) bool {
	if _, ok := paths.Path(key); ok {
		// prevent duplicate keys
		for _, confPathKey := range paths.ConfPathKeys {
			if confPathKey==key {
				return true
			}
		}
		// set the key as a conf path
		paths.ConfPathKeys = append(paths.ConfPathKeys, key)
		return true
	}
	// key doesn't match any paths
	return false
}

// Match a subpath in each of the ConfPaths, and return an ordered array
func (paths *Paths) GetConfSubPaths(subPath string) []string {
	confPaths := []string{}

	for _, confPathkey := range paths.ConfPathKeys {
		if confPath, ok := paths.Path(confPathkey); ok {
			confPaths = append(confPaths, path.Join(confPath, subPath))
		}
	}

	return confPaths
}

// Check to see if a file/folder path exists
// @NOTE it's best not to check if a path exists.
func (paths *Paths) CheckFileExists(filePath string) bool {
	_, err := os.Stat(filePath);
	return err==nil
}