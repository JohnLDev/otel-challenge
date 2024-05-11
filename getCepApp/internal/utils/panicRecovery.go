package utils

import "fmt"

func PanicRecovery() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from panic")
	}
}
