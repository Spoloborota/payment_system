package main

import (
	"fmt"
	"github.com/spf13/viper"
	"payment-system/internal"
)

func main() {
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file not found: %s \n", err))
	}

	internal.StartSystem()
}
