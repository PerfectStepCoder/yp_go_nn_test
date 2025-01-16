package main

import (
	"fmt"
	"os"
)

type Settings struct {
	ServiceHostOne string
	ServicePortOne string
	ServiceHostTwo string
	ServicePortTwo string
}

func NewSettings() (*Settings, error) {

	serviceHostOne := os.Getenv("SERVICE_HOST_ONE")
	servicePortOne := os.Getenv("SERVICE_PORT_ONE")
	serviceHostTwo := os.Getenv("SERVICE_HOST_TWO")
	servicePortTwo := os.Getenv("SERVICE_PORT_TWO")

	return &Settings{
		ServiceHostOne: serviceHostOne,
		ServicePortOne: servicePortOne,
		ServiceHostTwo: serviceHostTwo,
		ServicePortTwo: servicePortTwo,
	}, nil
}

func (s Settings) String() string {
	result := fmt.Sprintf("Settings:\n\tServiceHostOne:%s\n\tServicePortOne:%s\n\tServiceHostTwo:%s\n\tServicePortTwo:%s\n",
		s.ServiceHostOne, s.ServicePortOne, s.ServiceHostTwo, s.ServicePortTwo)
	return result
}
