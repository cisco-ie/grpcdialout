# Cisco MDT gRPC Dialout Collector

The Cisco MDT (Model Driven Telemetry) gRPC Dialout Collector allows the collection of Cisco's MDT Streaming Telemetry through the gRPC interface on certain IOS XR, NX-OS, and IOS XE devices. The output of the collector is to either a file or to kafka. This collector can only handle the K/V format, not Compact GPB. This collector also does not support TLS at this time.

The intent of this collector is for testing MDT and getting started with it, this collector does not take scale into account and might fail in full production environments.

## Getting Started

### Config

The collector requires a JSON config file called config.json. Below is an example of the config.

```json
{
	"kafka": {
		"brokers": ["localhost:9092"],
		"topic": "TelemetryTest"
	},
	"raw": false,
	"dump": true,
	"filename": "telemetry.txt"
}
```

The keys `raw` and `dump` are the only required fields. Raw sends outputs the data in the ProtoBuf format rather than running it through the de-serializer. 

If `dump` is true, you can specify the filename if you want it different from the default of `telemetry.txt`. 

If you want to send the data to kafka you will need to speicfy both brokers and topic.

### Getting the tool

#### Using Go

`go get github.com/skkumaravel/grpcdialout`

Then `go build` or `go install`

#### Download the release

Use github to download the release

## Issues

If you run into problems, please use Github Issues and Pull Requests
