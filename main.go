package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MessageData struct {
	I1   float64 `json:"i1"`
	I2   float64 `json:"i2"`
	I3   float64 `json:"i3"`
	U1   float64 `json:"u1"`
	U2   float64 `json:"u2"`
	U3   float64 `json:"u3"`
	Time string  `json:"time"`
}

func publishMessage(client mqtt.Client, topic string, data MessageData, clientID string) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Fehler beim Umwandeln der Daten in JSON:", err)
		return
	}

	token := client.Publish(topic, 2, true, jsonData)
	token.Wait()

	coloredClientID := fmt.Sprintf("\x1b[32m%s\x1b[0m", clientID)
	fmt.Printf("%s: Published message: %s to topic: %s\n", coloredClientID, jsonData, topic)
}

func generateRandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func roundToTwoDecimals(value float64) float64 {
	return float64(int(value*100)) / 100.0
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	opts := client.OptionsReader()
	clientID := opts.ClientID()
	coloredClientID := fmt.Sprintf("\x1b[32m%s\x1b[0m", clientID)
	fmt.Printf("%s: Connected\n", coloredClientID)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func startClient(clientID string, topic string, broker string, username string, password string, wg *sync.WaitGroup) {
	defer wg.Done()

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetUsername(username).
		SetPassword(password).
		SetConnectionLostHandler(connectLostHandler).
		SetOnConnectHandler(connectHandler).
		SetDefaultPublishHandler(messagePubHandler)

	opts.SetKeepAlive(2 * time.Second)
	opts.SetTLSConfig(&tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
	})

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		messageData := MessageData{
			I1:   roundToTwoDecimals(generateRandomFloat(0, 10000)),
			I2:   roundToTwoDecimals(generateRandomFloat(0, 10000)),
			I3:   roundToTwoDecimals(generateRandomFloat(0, 10000)),
			U1:   roundToTwoDecimals(generateRandomFloat(0, 10000)),
			U2:   roundToTwoDecimals(generateRandomFloat(0, 10000)),
			U3:   roundToTwoDecimals(generateRandomFloat(0, 10000)),
			Time: time.Now().Format("2006-01-02 15:04:05"),
		}

		publishMessage(client, topic, messageData, clientID)

		sleepDuration := time.Duration(rand.Intn(11)+1) * time.Second
		time.Sleep(sleepDuration)
	}
}

func main() {

	var showHelp bool
	var numClients int
	var broker string
	var topic string
	var username string
	var password string

	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.IntVar(&numClients, "clients", 0, "Number of MQTT clients to start")
	flag.StringVar(&broker, "broker", "", "Broker-address")
	flag.StringVar(&topic, "topic", "", "MQTT-topic")
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "user-password")

	flag.Parse()

	if showHelp || len(os.Args) == 1 {
		printHelp()
		return
	}

	if numClients <= 0 {
		fmt.Println("Error: The '-clients' argument is required and must be greater than 0.")
		flag.PrintDefaults()
		return
	}

	if broker == "" {
		fmt.Println("Error: The '-broker' argument is required.")
		flag.PrintDefaults()
		return
	}

	if topic == "" {
		fmt.Println("Error: The '-topic' argument is required.")
		flag.PrintDefaults()
		return
	}

	if username == "" {
		fmt.Println("Error: The '-username' argument is required.")
		flag.PrintDefaults()
		return
	}

	if password == "" {
		fmt.Println("Error: The '-password' argument is required.")
		flag.PrintDefaults()
		return
	}

	var wg sync.WaitGroup

	for i := 1; i <= numClients; i++ {
		wg.Add(1)
		go startClient(fmt.Sprintf("client%d", i), topic, broker, username, password, &wg)
	}

	wg.Wait()
}

func printHelp() {
	fmt.Println("")
	fmt.Println("MqttStorm - A tool for storming MQTT brokers with multiple clients")
	fmt.Println("\nUsage:")
	fmt.Println("  mqttstorm [options]")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}
