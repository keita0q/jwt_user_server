package main

import (
	"encoding/json"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"os"
	"github.com/keita0q/user_server/server"
	"github.com/keita0q/user_server/database/sequreDB/demoDB"
	"github.com/keita0q/user_server/database/applicationDatabase/demo"
	"github.com/keita0q/user_server/auth/jwtAuth"
)

func main() {
	tApplication := &cli.App{
		Name:    "user-server",
		Usage:   "user server with auth",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "config.json",
				Usage: "configuration file path",
			},
		},
		Action: func(aContext *cli.Context) error {
			tBytes, tError := ioutil.ReadFile(aContext.String("config"))
			if tError != nil {
				return tError
			}

			tConfig := Config{}
			if tError := json.Unmarshal(tBytes, &tConfig); tError != nil {
				return tError
			}
			tDB := demo.New()
			tSequreDB := demoDB.New()
			tAuth := jwtAuth.New(&jwtAuth.Config{
				DB: tSequreDB,
				PublicKeyPath: tConfig.PublicKeyPath,
				PrivateKeyPath: tConfig.PrivateKeyPath,
			})

			return server.Run(&server.Config{
				ContextPath: tConfig.ContextPath,
				Port:        tConfig.Port,
				ClientDir:   tConfig.ClientDir,
				Database:       tDB,
				SequreDB: tSequreDB,
				Auth: tAuth,
			})
		},
	}

	tApplication.Run(os.Args)
}

type Config struct {
	ContextPath    string   `json:"context_path"`
	Port           int      `json:"port"`
	ClientDir      string   `json:"client_dir"`
	PublicKeyPath  string   `json:"public_key_path"`
	PrivateKeyPath string   `json:"private_key_path"`
}