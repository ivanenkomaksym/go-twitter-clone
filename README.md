# Twitter clone in GO

Simple twitter clone with React frontend, Go backend and Server-Sent Events to support real-time refreshing. This example is based on [HTTP Server push using SSE (Server-Sent Events)](https://github.com/ThreeDotsLabs/watermill/tree/master/_examples/real-world-examples/server-sent-events)

![Alt text](docs/screen.gif?raw=true)

## Features

* Implements common messaging patterns such as publish-subscribe, request-reply, and event-driven architecture in Go.
* Utilizes popular database and  message broker technologies, including MySQL, Mongo and NATS.
* Using appsettings in Go applications.
* Includes Dockerfiles and Docker Compose configuration for containerizing the sample microservices.
* Supports basic Google OAuth2
* Runs database integration tests during CI using github workflow actions
* Runs end-to-end tests on full scale environment during CI using github workflow actions

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

### Integration tests in CI using github workflow actions
There's a dedicated workflow installed on repository that would do the following:
1. Spin up required services: MySQL and Mongo
2. Check out code
3. Set up Go environment
4. Build
5. Run all tests

## End-to-end tests

* Server supports two authentication methods:
1. Cookie-based authentication (visualized on the diagram above)
2. Authorization header that can be used by server application (a.k.a. daemon app) without a real user entering credentials
* In Google Cloud Console there's a service account configured that would let server application to authenticate. [Using OAuth 2.0 for Server to Server Applications](https://developers.google.com/identity/protocols/oauth2/service-account#httprest). It also requires a Key to be created, which is basically a JSON document that you download.

![Alt text](docs/authn.png?raw=true "Authentication")

* e2e tests rely on `GOOGLE_SERVICE_ACCOUNT_KEY` environment variable to generate JWT and then access token. You need to load content of private key JSON file into that environment variable, for example for PowerShell this could be:

```ps
$env:GOOGLE_SERVICE_ACCOUNT_KEY = Get-Content -Raw -Path "C:\Users\Maksym Ivanenko\Downloads\twitter-clone-438407-23XXXXXXXXX.json"
```

![Alt text](docs/privatekey.png?raw=true "Private key JSON")

* To run e2e test locally execute in `e2e_tests` folder:

```sh
npm test
```

* In order to run e2e tests during CI private key JSON has been loaded into github secret under the same name. Github action then forwards it as environment variable:
```yaml
- name: Run e2e tests
      env:
        GOOGLE_SERVICE_ACCOUNT_KEY: ${{ secrets.GOOGLE_SERVICE_ACCOUNT_KEY }}
```

* e2e test action does several things:
1. Checks out source code
2. Runs docker-compose to start full scale environment. It uses github secret for google service account key.
3. Installs NodeJS and dependencies in e2e_tests
4. Runs the test with `npm test`

![Alt text](docs/e2e_test.png?raw=true "End-to-end test")

# References
[HTTP Server push using SSE (Server-Sent Events)](https://github.com/ThreeDotsLabs/watermill/tree/master/_examples/real-world-examples/server-sent-events)

[Google OAuth 2.0 and Golang](https://medium.com/@_RR/google-oauth-2-0-and-golang-4cc299f8c1ed)

[Using OAuth 2.0 for Server to Server Applications](https://developers.google.com/identity/protocols/oauth2/service-account#httprest)