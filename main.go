package main

import (
	"log"
	"os"
	"text/template"
)

type Config struct {
	WanMac           string
	LanMac           string
	SSHAuthorizedKey string
}

func main() {
	tmpl, err := template.ParseFiles("templates/vyos-postconfig-bootup.script.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	wanMac := os.Getenv("WAN_MAC")
	lanMac := os.Getenv("LAN_MAC")
	sshAuthorizedKey := os.Getenv("SSH_AUTHORIZED_KEY")

	if wanMac == "" || lanMac == "" || sshAuthorizedKey == "" {
		panic("Environment variable not set. Make all three have a value: WAN_MAC, LAN_MAC, SSH_AUTHORIZED_KEY")
	}

	data := Config{
		WanMac:           wanMac,
		LanMac:           lanMac,
		SSHAuthorizedKey: sshAuthorizedKey,
	}

	file, err := os.Create("vyos-postconfig-bootup.script")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		log.Fatal(err)
	}
}
