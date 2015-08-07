package main

import (
// 	"os"
// 	"io"
//
// 	"github.com/libgit2/git2go"
)

func (operation *Operation_Init) Init_Git_Run(flags []string) (bool, map[string]string) {



	return true, map[string]string{
		".coach/CREATEDFROM.md":  `THIS PROJECT WAS CREATED FROM GIT`,
	}
}
