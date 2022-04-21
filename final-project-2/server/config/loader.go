package configuration

import (
	"fmt"
	"os"
)

// Return default PORT number
func GetPort(port string) string {
	if port == "" {
		return ":8000"
	} else {
		return fmt.Sprintf(":%v", os.Getenv("PORT"))
	}
}
