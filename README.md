# docker-compose example on instrumenting the go application for prometheus
This Project is basic example on how to write prometheus metric generator/collector using golang and docker-compose.
This gives starting point on how to instrument the go application for prometheus. Also, it will give opportunity in trying various promql queries on data controlled by user(main.go). It helps in giving more insights about the query.

## Pre-requisites

Ensure you install the latest version of make, docker and [docker-compose](https://docs.docker.com/compose/install/) on your Docker host machine. This has been tested with Docker for Mac.

# Quick Start

The code present is in main.go file
This has three collectors registered:

```
prometheus.Register(version.NewCollector(collector))
prometheus.Register(NewSysStatCollector(collector))
prometheus.Register(taskCounterVec)
```

To build and run the docker-compose :

1. (optional) add your code to main.go file
2. execute `make image` from top of repo
3. cd monitor
4. docker-compose up -d

Your docker-compose will be up and running.
You can check prom metric using `http://localhost:9090` in your browser.
`http://localhost:9090/targets` shows status of monitored targets as seen from prometheus.

## Reference and disclaimer:
For more information on Prometheus instrumentation :   [Prometheus](https://prometheus.io/docs/guides/go-application/)
Disclaimer: this setup is not in anyway secured. So use it only for inspiration.