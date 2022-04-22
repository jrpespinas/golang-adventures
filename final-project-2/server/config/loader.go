package configuration

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
}

// Return default PORT number
func GetPort(port string) string {
	if port == "" {
		return ":3333"
	} else {
		return fmt.Sprintf(":%v", os.Getenv("PORT"))
	}
}
