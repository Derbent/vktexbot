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
		log.Print("VK_TOKEN env variable does not exist!")
		return
	}

	groupID, exists := os.LookupEnv("VK_GROUP_ID")
	if exists == false {
		log.Print("VK_GROUP_ID env variable does not exist!")
		return
	}

	intGroupID, err := strconv.Atoi(groupID)
	if err != nil {
		log.Print("Bad VK_GROUP_ID")
		return
	}

	vk, err := transport.NewVK(token, intGroupID)
	if err != nil {
		log.Print(err)
		return
	}
	transport.Run(vk)
	log.Print("VK server stopped")
}

func main() {
	RunVKServer()
}
