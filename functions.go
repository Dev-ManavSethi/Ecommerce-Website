package main

import "log"

func FatalOnError(message string, err error) {
	if err != nil {
		log.Println(message)
		log.Println(err)

	}
}
