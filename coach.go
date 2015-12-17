package main

import "github.com/james-nesbitt/v2/log"
	// "conf"
	// "node"
	// "client"

import	"fmt"


var Locallog *CoachLog
  // conf *Conf
  // client *Client
  // nodes *Nodes


func init() {

	// conf = getConf()
	// client = getClient()
	// nodes = getNodes()

}

func main() {


	log = getLog()	

  fmt.Println("LOG:", log)
  // fmt.Println("CONF:", conf)
  // fmt.Println("CLIENT:", client)
  // fmt.Println("NODES:", nodes)
  
}
