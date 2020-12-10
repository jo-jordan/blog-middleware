package common

import (
	"log"
)

func ErrorBus(err error) {
	if err != nil {
		log.Printf("err: %s", err)
	}
}
