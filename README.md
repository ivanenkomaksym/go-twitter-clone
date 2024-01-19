# Twitter clone in GO

Simple twitter clone with React frontend, Go backend and Server-Sent Events to support real-time refreshing. This example is based on [HTTP Server push using SSE (Server-Sent Events)](https://github.com/ThreeDotsLabs/watermill/tree/master/_examples/real-world-examples/server-sent-events)

![Alt text](docs/screen.gif?raw=true)

## Features

* Implements common messaging patterns such as publish-subscribe, request-reply, and event-driven architecture in Go.
* Utilizes popular database and  message broker technologies, including MySQL, Mongo and NATS.
* Using appsettings in Go applications.
* Includes Dockerfiles and Docker Compose configuration for containerizing the sample microservices.

![Alt text](docs/architecture.png?raw=true "Application architecture")

# Documentation

## Settings
Application settings can be configured in `internal/config/appsettings.json` file. Possible configurations include mode of the application, API server configuration, Tweets and Feeds database connections, NATS messaging connection.

# Run in a Docker Container

To run the service using Docker Compose, use the following command:

```
docker-compose up
```

This command launches the service along with its dependencies defined in the docker-compose.yml file.

# References
[HTTP Server push using SSE (Server-Sent Events)](https://github.com/ThreeDotsLabs/watermill/tree/master/_examples/real-world-examples/server-sent-events)