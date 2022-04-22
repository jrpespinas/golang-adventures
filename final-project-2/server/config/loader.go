package configuration

import (
	"fmt"
	"os"
)

// Return default PORT number
func GetPort(port string) string {
	if port == "" {
		return ":3333"
	} else {
		return fmt.Sprintf(":%v", os.Getenv("PORT"))
	}
}
