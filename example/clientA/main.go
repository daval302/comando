package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Conf struct {
	Name string
	ID   string
}

var conf Conf

func main() {

	// Load configuration json file
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Load configuration file")
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	json.Unmarshal(byteValue, &conf)

	// get 2 bytes json as ID
	clientID := []byte(conf.ID)

	//fmt.Printf("client : %#x\n", clientID)
	fmt.Print(clientID)
}
