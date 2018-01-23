package engine

import "log"

func LogHandler(event *Event) (*Event, error) {
	log.Println(event)
	return nil,nil
}
