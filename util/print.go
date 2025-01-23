package util

import (
	"fmt"
)

func Eprint(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func EMsg(err error, msg string) {
	if err != nil {
		fmt.Printf("error: %s --> %s\n", msg, err)
	}
}
