# Twitter clone in GO

![Alt text](docs/ui.png?raw=true)

Simple twitter clone with React frontend, Go backend and Server-Sent Events to support real-time refreshing. This example is based on [HTTP Server push using SSE (Server-Sent Events)](https://github.com/ThreeDotsLabs/watermill/tree/master/_examples/real-world-examples/server-sent-events). Supports multi-environment setup, containerization, CI/CD, cloud deployment.

![Alt text](docs/screen.gif?raw=true)

## Features

* Implements common messaging patterns such as publish-subscribe, request-reply, and event-driven architecture in Go.
* Complete CI/CD pipeline: builds both backend and frontend services, runs unit and database integration tests, executes end-to-end tests in a full-scale environment, and deploys to Google Cloud Run.
* Support FireStore and Pub/Sub in Cloud setup.
* Hosted on Google Cloud DNS.
* Utilizes popular database and  message broker technologies, including MySQL, Mongo and NATS.
* Using appsettings in Go applications.
* Includes Dockerfiles and Docker Compose configuration for containerizing the sample microservices.
* Supports basic Google OAuth2.
* Runs database integration tests during CI using github workflow actions.
* Includes common project structure for frontend projects.
* Runs frontend unit tests.

![Alt text](docs/architecture.png?raw=true "Application architecture")

![Alt text](docs/authn.png?raw=true "Authentication")

# Documentation

## Settings
Application settings can be configured in `internal/config/appsettings.json` file. Possible configurations include mode of the application, API server configuration, Tweets and Feeds database connections, NATS messaging connection.

## Run cloud mode locally
First set `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to service account key, for example for PowerShell this could be:

```ps
$env:GOOGLE_APPLICATION_CREDENTIALS = "D:\Downloads\twitter-clone-438407-4b74b07d276e.json"
```

Then in `server` folder run
```
go run cmd\main.go
```

# Multi-Environment Configuration for Deployment

This project supports multiple deployment environments: in-memory, persistent, and cloud using Docker Compose. CI/CD workflow supports cloud deployment environment by default.

![Alt text](docs/envs.png?raw=true "Multi-Environment Configuration")

This setup provides a robust and flexible way to handle multiple environments while ensuring security and consistency across development and production workflows.

Below is a detailed explanation of the setup:

## Environment-Specific Configuration Files
Environment variables are managed using `.env` files tailored for each deployment mode:
- **`.env`**: Default environment for in-memory mode, suitable for local development.
- **`.env.local`**: Persistent mode, configured for local development with a database and messaging connections.
- **`.env.cloud`**: Cloud mode, used for deployment in cloud environment.

## Docker Compose Setup
The `compose.yaml` references several environment variables that can be dynamically loaded from the specified `.env` file using the `--env-file` option.

Example command:
```bash
docker-compose --env-file .env.local -f compose.yaml up -d --build
```

## CI/CD Pipeline Integration
The `ci-cd.yml` file integrates the configuration into a CI/CD workflow, using secrets stored securely in a GitHub Actions. Environment-specific secrets are injected during pipeline execution:
- `DOCKER_USER` and `DOCKER_PASSWORD` for Docker image management.
- `PROJECT_ID`, `OAUTH2_CLIENT_ID`, `OAUTH2_CLIENT_SECRET` `GOOGLE_SERVICE_ACCOUNT_KEY` for cloud configuration.

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

### Router tests

`mockgen` can generate necessary mocks that can be used for testing:
```
mockgen -source .\internal\api\publisher.go -destination .\internal\__mocks__\api\mock_publisher.go -package=mocks
```

Router test can be run with
```
go test .\internal\api\router_test.go
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

# Frontend

## Project structure

```
src/
├── __mocks__/
│   └── configMock.js                 # Mocks common.js
│
├── __tests__/
│   ├── apiHandlers.tests.js          # Tests for API request functions
│   └── eventSourceHandlers.test.js   # Tests for even source request functions
│
├── api/
│   ├── apiHandlers.js                # Contains all API request functions
│   └── eventSourceHandlers.js        # Contains all event source request functions
│
├── components/
│   ├── auth/
│   │   └── AuthContext.tsx           # Contains AuthContext, checkAuth function
│   │
│   ├── pages/
│   │   ├── Callback.js               # Callback page for authentication redirection
│   │   ├── Logic.js                  # Login page
│   │   ├── Main.js                   # Main page combining TagList, TweetList, InputForm
│   │   ├── Nav.js                    # Main top navigation page
│   │   └── Profile.js                # Profile page showing user info from AuthContext
│   │
│   ├── tweet/
│   │   ├── InputForm.js              # Form for creating a new tweet
│   │   ├── TagList.js                # Displays list of hashtags
│   │   ├── TweetCard.jsx             # Represents a single tweet item in TweetList
│   │   └── TweetList.js              # Displays a list of TweetCard components
│   │
├── models/
│   └── user.js                       # User model
│
├── styles/
│   ├── pages/                        # Styles specific to page components
│   └── tweet/                        # Styles specific to tweet-related components
│
├── App.js                            # Main app entry point
├── common.js                         # Loads environemnt variables and defines constants
├── env.json                          # Environment configuration
└── index.js                          # Entry point for rendering the React app
```

# References
[HTTP Server push using SSE (Server-Sent Events)](https://github.com/ThreeDotsLabs/watermill/tree/master/_examples/real-world-examples/server-sent-events)

[Google OAuth 2.0 and Golang](https://medium.com/@_RR/google-oauth-2-0-and-golang-4cc299f8c1ed)

[Using OAuth 2.0 for Server to Server Applications](https://developers.google.com/identity/protocols/oauth2/service-account#httprest)

[Set up a global external Application Load Balancer with Cloud Run, App Engine, or Cloud Run functions](https://cloud.google.com/load-balancing/docs/https/setup-global-ext-https-serverless)

[Tutorial: Set up a domain by using Cloud DNS](https://cloud.google.com/dns/docs/tutorials/create-domain-tutorial#set-up-domain)