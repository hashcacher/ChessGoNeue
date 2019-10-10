package core
import (
	"os"
	"fmt"
)

func Debug(msg string) {
	if os.Getenv("DEBUG") == "true" {
		fmt.Println(msg)
	}
}
