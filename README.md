# MqttStorm

MqttStorm is a tool for storming MQTT brokers with multiple clients.

## Usage

```
mqttstorm [options]
```

## Options

- `-help`: Show help message
- `-clients`: Number of MQTT clients to start
- `-broker`: Broker address
- `-topic`: MQTT topic
- `-username`: Username
- `-password`: User password

## How to Use

1. Clone the repository:

```bash
git clone https://github.com/schneider-marco/mqttstorm
```

2. Compile the code:

```bash
cd mqttstorm
go build
```

3. Run the executable:

```bash
./mqttstorm [options]
```

Make sure to replace `[options]` with the appropriate values for your MQTT setup.

## License

This project is licensed under the [MIT License](LICENSE).
