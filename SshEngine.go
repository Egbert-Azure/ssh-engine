package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

func main() {
	// Setup logging
	file, err := os.OpenFile("engine.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)

	// log.Println("Running on " + runtime.GOOS)

	// Read configuration
	configuration := readConfiguration()

	server := fmt.Sprintf("%s:%s", configuration.Host, configuration.Port)

	// Setup the client configuration
	sshConfig := getSshConfig(configuration)

	// Start the connection
	client, err := ssh.Dial("tcp", server, sshConfig)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	// Start a session
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// StdinPipe for commands
	stdin, _ := session.StdinPipe()

	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %s", err)
	}

	// Run the supplied command first
	fmt.Fprintf(stdin, "%s\n", configuration.RemoteCommand)

	// Accepting commands
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		log.Println("Input: " + scanner.Text())
		fmt.Fprintf(stdin, "%s\n", scanner.Text())
		if scanner.Text() == "quit" {
			log.Println("Quit sent")
			break
		}
	}
}

func getSshConfig(configuration Configurations) *ssh.ClientConfig {
	key, err := getKeyFile(configuration.PrivateKeyFile)
	if err != nil {
		panic(err)
	}

	sshConfig := &ssh.ClientConfig{
		User: configuration.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	return sshConfig
}

func getKeyFile(file string) (key ssh.Signer, err error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}

func readConfiguration() Configurations {
	viper.SetConfigName("engine")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	var configuration Configurations
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return configuration
}

type Configurations struct {
	User           string
	PrivateKeyFile string
	Host           string
	Port           string
	RemoteCommand  string
}
