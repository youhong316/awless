package main

import (
	"log"

	"github.com/spf13/viper"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/cloud/azure"
	"github.com/wallix/awless/cmd"
	"github.com/wallix/awless/config"
)

func main() {
	config.InitAwlessEnv()

	sess, err := aws.InitSession(viper.GetString("region"))
	if err != nil {
		log.Fatal(err)
	}

	aws.InitServices(sess)

	if sess, err := azure.InitSession(); err == nil {
		azure.InitServices(sess)
	}

	cmd.InitCli()
	cmd.ExecuteRoot()
}
