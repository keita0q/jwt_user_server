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
	"github.com/keita0q/user_server/mail/smtp"
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

			tMail := smtp.New(&smtp.Config{
				Address: tConfig.MailConfig.Address,
				NeedAuth: tConfig.MailConfig.NeedAuth,
				Identity: tConfig.MailConfig.Identity,
				UserName: tConfig.MailConfig.Username,
				Password: tConfig.MailConfig.Password,
				Host: tConfig.MailConfig.Host,
			})
			return server.Run(&server.Config{
				ContextPath: tConfig.ContextPath,
				Port:        tConfig.Port,
				ClientDir:   tConfig.ClientDir,
				Database:       tDB,
				SequreDB: tSequreDB,
				Auth: tAuth,
				Mail: tMail,
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
	MailConfig     MailConfig `json:"mail_config"`
}

type MailConfig struct {
	Address  string `json:"address"`
	NeedAuth bool   `json:"need_auth"`
	Identity string `json:"identity"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}