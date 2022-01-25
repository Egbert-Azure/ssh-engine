package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

func main() {
	// log.Println("Running on " + runtime.GOOS)

	// Read configuration
	configuration := readConfiguration()
	debugLogging := false

	// Setup logging if a log file name was passed in
	if configuration.LogFileName != "" {
		file, err := os.OpenFile("engine.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		log.SetOutput(file)

		debugLogging = true
	}

	server := fmt.Sprintf("%s:%s", configuration.Host, configuration.Port)

	// Setup the client configuration
	sshConfig := getSshConfig(configuration)

	// Start the connection
	client, err := ssh.Dial("tcp", server, sshConfig)
	if err != nil {
		fmt.Printf("Could not connect to ssh (failed to dial). Error is: %s\n", err)
		os.Exit(1)
	}

	// Start a session
	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create ssh session. Error is: %s\n", err)
		os.Exit(1)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// StdinPipe for commands
	stdin, _ := session.StdinPipe()

	// Start remote shell
	if err := session.Shell(); err != nil {
		fmt.Printf("Failed to start shell. Error is: %s\n", err)
		os.Exit(1)
	}

	// Run the supplied command first
	fmt.Fprintf(stdin, "%s\n", configuration.RemoteCommand)

	// Accepting commands
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		if debugLogging {
			log.Println("Input: " + scanner.Text())
		}

		fmt.Fprintf(stdin, "%s\n", scanner.Text())
		if scanner.Text() == "quit" {
			if debugLogging {
				log.Println("Quit sent")
			}

			break
		}
	}
}

func getSshConfig(configuration Configurations) *ssh.ClientConfig {
	key, err := getKeyFile(configuration.PrivateKeyFile)
	if err != nil {
		fmt.Printf("Could not read privateKeyFile at %s\n", configuration.PrivateKeyFile)
		os.Exit(1)
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
		fmt.Printf("Error reading the key file. Error is: %s\n", err)
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		fmt.Printf("Error parsing the private key file. Is this a valid private key? Error is: %s\n", err)
		return
	}
	return
}

func readConfiguration() Configurations {
	if _, err := os.Stat("engine.yml"); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		fmt.Println("The file 'engine.yml' could not be found in the current directory")
		os.Exit(1)
	}

	// Set the file name of the configurations file
	viper.SetConfigName("engine")
	viper.SetConfigType("yml")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Read the Configuration
	var configuration Configurations
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("No such config file")
			os.Exit(1)
		} else {
			// Config file was found but another error was produced
			fmt.Printf("Error reading the engine.yml file, %s", err)
			os.Exit(1)
		}

	}

	// Set undefined variables
	viper.SetDefault("logFileName", "")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode the engine.yml file, %v", err)
		os.Exit(1)
	}

	return configuration
}

type Configurations struct {
	User           string
	PrivateKeyFile string
	Host           string
	Port           string
	RemoteCommand  string
	LogFileName    string
}
