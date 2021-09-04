package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ZashX/vktexbot/transport"
)

func RunVKServer() {
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

	vk := transport.NewVK(token, intGroupID)
	transport.Run(vk)
}

func main() {
	RunVKServer()
	fmt.Print("RUN")
}
