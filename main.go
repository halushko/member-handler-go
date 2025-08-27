package member_handler_go

import (
	"log"

	"github.com/halushko/core-go/logger"
	"github.com/halushko/member-handler-go/database"
	"github.com/halushko/member-handler-go/handlers"
)

//goland:noinspection GoUnusedFunction
func main() {
	logFile := logger.SoftPrepareLogFile("member_handler_go.log")

	log.Println("Starting member_handler")

	err := database.Init()
	if err != nil {
		log.Printf("[ERROR] Start of member_handler failed: %v", err)
		return
	}

	go handlers.StartMemberJoinedListener()

	defer logger.SoftLogClose(logFile)

	select {}
}
