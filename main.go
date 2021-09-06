package main

import (
	"log"
	"os"
	"strconv"

	"github.com/ZashX/vktexbot/transport"
)

func RunVKServer() {
	log.Print("Start VK server")

	token, exists := os.LookupEnv("VK_TOKEN")
	if exists == false {
		panic("VK_TOKEN env variable does not exist!")
	}

	groupID, exists := os.LookupEnv("VK_GROUP_ID")
	if exists == false {
		panic("VK_GROUP_ID env variable does not exist!")
	}

	intGroupID, err := strconv.Atoi(groupID)
	if err != nil {
		panic("Bad VK_GROUP_ID")
	}

	vk, err := transport.NewVK(token, intGroupID)
	if err != nil {
		panic(err)
	}
	transport.Run(vk)
	log.Print("VK server stopped")
}

func main() {
	RunVKServer()
}
