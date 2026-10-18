package main

import (
	"log"
	"os"
	"strconv"
	"text/template"
)

type MDNSConfig struct {
	Enabled       bool
	InterfaceMDNS string // ethX
	Hostname      string
}

type BaseConfig struct {
	WanMac           string
	LanMac           string
	SSHAuthorizedKey string
}

func main() {
	wanMac := os.Getenv("WAN_MAC")
	lanMac := os.Getenv("LAN_MAC")
	sshAuthorizedKey := os.Getenv("SSH_AUTHORIZED_KEY")
	mdnsEnable := os.Getenv("ENABLE_MDNS")

	if wanMac == "" || lanMac == "" || sshAuthorizedKey == "" {
		panic("Environment variable not set. Make all three have a value: WAN_MAC, LAN_MAC, SSH_AUTHORIZED_KEY")
	}

	CreateBaseConfigFile(
		"templates/vyos-postconfig-bootup.script.tmpl",
		"scripts/vyos-postconfig-bootup.script",
		BaseConfig{
			WanMac:           wanMac,
			LanMac:           lanMac,
			SSHAuthorizedKey: sshAuthorizedKey,
		},
	)

	CreateMDNSConfigFile(
		"scripts/vyos-postconfig-bootup.script",
		"scripts/vyos-postconfig-bootup.script",
		MDNSConfig{
			Hostname:      "",
			InterfaceMDNS: "",
			Enabled:       ParseEnvVarBool(mdnsEnable),
		},
	)
}

func CreateBaseConfigFile(templatePath, destPath string, config BaseConfig) {
	tmpl, err := template.ParseFiles()
	if err != nil {
		log.Fatal(err)
	}

	data := BaseConfig{}

	file, err := os.Create(destPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		log.Fatal(err)
	}
}

func CreateMDNSConfigFile(templatePath, destPath string, config MDNSConfig) {
	if !config.Enabled {
		return
	}
}

func ParseEnvVarBool(envVar string) bool {
	if envVar != "true" || envVar != "false" {
		log.Panicf("Invalid values. Only 'true' or 'false' are available")
	}

	b, err := strconv.ParseBool(envVar)
	if err != nil {
		log.Panicf("")
	}

	return b
}
