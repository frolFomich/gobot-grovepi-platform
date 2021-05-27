# gobot-grovepi-platform
[gobot.io](https://gobot.io/) based [GrovePi](https://www.dexterindustries.com/GrovePi/get-started-with-the-grovepi/) platform for [IoT](https://en.wikipedia.org/wiki/Internet_of_things) experiments

### How to build

#### CLI

To build this application Go environment should be setup
1. Clone this repo to go/src dir
1. Execute following command in root of cloned repo directory
> export GOOS=linux GOARCH=arm64; go build -o gobot-grovepi-platform cmd/grovepi/main.go

#### Docker image

To build docker image on macOS execute following command
> docker buildx build --platform linux/arm --rm --tag <your-docker-registry>/gobot-grovepi-platform:latest --push -f deploy/Dockerfile .

### How to run

#### CLI

1. Copy built command and config file `config/app.yaml` to RaspberryPi
1. Execute following command
> gobot-grovepi-platform -c app.yaml
1. You may access the robeaux React.js interface with Gobot by navigating to http://localhost:3000/index.html.

#### Docker image

> docker run --rm --privileged -p 3000:3000 <your-docker-registry>/gobot-grovepi-platform:latest

### Disclaimer

Working with such hardware like RaspberryPi/GrovePi/other may be dangerous for inexperienced people.
Author isn't liable for any potential damages caused by this software. Use it with caution at your own risk 
