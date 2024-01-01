# Kindle Smart Home Info Display

Fetch data such as weather forecast, train departures or other smart home data and print it nicely.

In contrast to other similar projects, this completely runs on the Kindle itself and requires no server component.

For more information about jailbreaking and building a physical frame, see the [accompanying blog post](https://blog.simon-dreher.de/kindle-info-dashboard.html).

# Build

For trying locally use `go run main.go` and open the resulting `example.png`.

To build for Kindle, you need to use golang 1.17 (newer versions don't run on Kindles old kernel).
You can use a container to do this:
```
podman run --rm -v "${PWD}":/app -it docker.io/golang:1.17.13 bash -c "cd /app && env GOOS=linux GOARCH=arm GOARM=5 go build"
```
