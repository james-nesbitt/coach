package main

import (
	"os"
	"io"
	"path"
)

func (operation *Operation_Init) Init_User_Run(flags []string) (bool, map[string]string) {
operation.log.DebugObject( LOG_SEVERITY_MESSAGE, "FLAGS:", flags)

	template := flags[0]
	flags = flags[1:]

	templatePath, ok := operation.conf.Path("usertemplates")
	if !ok {
		operation.log.Error("COACH has no directive for a user template path")
		return false, map[string]string{}
	}
	targetPath := path.Join(templatePath , template )

	if _, err := os.Stat( targetPath ); err!=nil {
		operation.log.Error("Invalid template path suggested for new project init : ["+template+"] expected path ["+targetPath+"] => "+err.Error())
		return false, map[string]string{}
	}

	operation.log.Message("Copying new project init from template : ["+template+"] path ["+targetPath+"]")
	if err := CopyDir(targetPath, operation.root); err!=nil {
		operation.log.Error("Falied copying new project init from template : ["+template+"] expected path ["+targetPath+"] => "+err.Error())
		return false, map[string]string{}
	}

	return true, map[string]string{
		".coach/CREATEDFROM.md":  `THIS PROJECT WAS CREATED FROM A User Template`,
	}
}

func CopyFile(source string, dest string) (err error) {
		sourcefile, err := os.Open(source)
		if err != nil {
				return err
		}

		defer sourcefile.Close()

		destfile, err := os.Create(dest)
		if err != nil {
				return err
		}

		defer destfile.Close()

		_, err = io.Copy(destfile, sourcefile)
		if err == nil {
				sourceinfo, err := os.Stat(source)
				if err != nil {
						err = os.Chmod(dest, sourceinfo.Mode())
				}

		}

		return
}

func CopyDir(source string, dest string) (err error) {

		// get properties of source dir
		sourceinfo, err := os.Stat(source)
		if err != nil {
				return err
		}

		// create dest dir

		err = os.MkdirAll(dest, sourceinfo.Mode())
		if err != nil {
				return err
		}

		directory, _ := os.Open(source)

		objects, err := directory.Readdir(-1)

		for _, obj := range objects {

				sourcefilepointer := source + "/" + obj.Name()

				destinationfilepointer := dest + "/" + obj.Name()


				if obj.IsDir() {
						// create sub-directories - recursively
						err = CopyDir(sourcefilepointer, destinationfilepointer)
						if err != nil {
// 								fmt.Println(err)
						}
				} else {
						// perform copy
						err = CopyFile(sourcefilepointer, destinationfilepointer)
						if err != nil {
// 								fmt.Println(err)
						}
				}

		}
		return
}
