# Twitter clone in GO

Simple twitter clone with React frontend, Go backend and Server-Sent Events to support real-time refreshing. This example is based on [HTTP Server push using SSE (Server-Sent Events)](https://github.com/ThreeDotsLabs/watermill/tree/master/_examples/real-world-examples/server-sent-events)

![Alt text](docs/screen.gif?raw=true)

## Features

* Implements common messaging patterns such as publish-subscribe, request-reply, and event-driven architecture in Go.
* Utilizes popular database and  message broker technologies, including MySQL, Mongo and NATS.
* Using appsettings in Go applications.
* Includes Dockerfiles and Docker Compose configuration for containerizing the sample microservices.
* Supports basic Google OAuth2

![Alt text](docs/architecture.png?raw=true "Application architecture")

![Alt text](docs/authn.png?raw=true "Authentication")

# Documentation

## Settings
Application settings can be configured in `internal/config/appsettings.json` file. Possible configurations include mode of the application, API server configuration, Tweets and Feeds database connections, NATS messaging connection.

# Run in a Docker Container

To run the service using Docker Compose, use the following command:

```
docker-compose up
```

This command launches the service along with its dependencies defined in the docker-compose.yml file.

# Testing

To run all existing tests recursively execute in the root folder:

```
go test ./... -v
```

## Integration tests

### Feed tests

First make sure you have Mongo instance running, e.g. using docker:
```
docker run -d -p 27017:27017 --name mongo mongo:latest
```

Feed integation test will connect to `Tests_FeedsDb` Mongo database. To run tests execute in root folder:

```
go test .\internal\repositories\feed\persistent_feedrepository_test.go -v
```

### Tweet tests

First make sure you have MySQL instance running, e.g. using docker:
```
docker run -d --name twitter-mysql-test -e MYSQL_ROOT_PASSWORD=my-secret-pw -e MYSQL_USER=myuser -e MYSQL_PASSWORD=mypassword -e MYSQL_DATABASE=Tests_TweetsDb -p 3306:3306 mysql:latest
```

Tweet integration test will connect to `Tests_TweetsDb` MySql database. To run tests execute in root folder:

```
go test .\internal\repositories\tweet\persistent_tweetrepository_test.go -v
```

# References
[HTTP Server push using SSE (Server-Sent Events)](https://github.com/ThreeDotsLabs/watermill/tree/master/_examples/real-world-examples/server-sent-events)

[Google OAuth 2.0 and Golang](https://medium.com/@_RR/google-oauth-2-0-and-golang-4cc299f8c1ed)