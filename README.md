# gobot-grovepi-platform
[gobot.io](https://gobot.io/) based [GrovePi](https://www.dexterindustries.com/GrovePi/get-started-with-the-grovepi/) platform for [IoT](https://en.wikipedia.org/wiki/Internet_of_things) experiments

### How to build

To build this application Go environment should be setup
1. Clone this repo to go/src dir
1. Execute following command in root of cloned repo directory
> export GOOS=linux GOARCH=arm64; go build -o gobot-grovepi-platform cmd/grovepi/main.go

### How to run

1. Copy built command and config file `config/app.yaml` to RaspberryPi
1. Execute following command
> main -f app.yaml
1. You may access the robeaux React.js interface with Gobot by navigating to http://localhost:3000/index.html.

### Disclaimer

Working with such hardware like RaspberryPi/GrovePi/other may be dangerous for inexperienced people.
Author isn't liable for any potential damages caused by this software. Use it with caution at your own risk 
