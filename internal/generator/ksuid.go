package generator

import "github.com/segmentio/ksuid"

func NewKSUID() string {
	return ksuid.New().String()
}
