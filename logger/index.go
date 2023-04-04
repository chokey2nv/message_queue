package logger

import "fmt"

func Log(err error) {
	//perform logging here or error notification

	fmt.Println(err.Error())
}
