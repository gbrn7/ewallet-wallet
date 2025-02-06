package main

import (
	"ewallet-wallet/cmd"
	"ewallet-wallet/helpers"
)

func main() {
	// load config
	helpers.SetupConfig()

	// load log
	helpers.SetupLogger()

	// load db
	helpers.SetupMySQL()

	// run grpc
	// go cmd.ServeGRPC()

	// run http
	cmd.ServeHttp()

}
