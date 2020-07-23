# Patal Web Backend
Patal Web is repository for https://palembangdigital.org/

## Getting Started

To run the project localy, make sure minimum requirements are fulfilled.

- Go version 1.10 or higher
- PostgreSQL version 12
- Docker (optional) -- see [here](https://docs.docker.com/get-docker/).

### Running in Local Machine

1. Make sure Go is installed as global command (first time only)

2. Clone this project and go to the root project to install all dependencies (first time only)
    ```bash
    // clone the repository
    > git clone git@github.com:palembang-digital/web-backend.git

    // change directory to root project web-frontend folder
    > cd web-backend
    > rename config-example.yaml to config.yaml

    // install all the dependencies
    > make dep
    ```

3. While still in root project build and run the app
    ```bash
    // build and run the app
    > make build
    > ./web-backend

    // now go to http://localhost:8006/ in your browser to check the app.
    ```

### Running PostgreSQL from Docker

1. Make sure Docker and Docker Compose is installed

2. Run `docker-compose up -d`, it will start the PostgreSQL DB as docker container.

3. Build and run the web-backend app as described on the previous section.

## API Documentation

We use [swag](https://github.com/swaggo/swag) to generate necearry Swagger files for API documentation. Everytime we run `make build`, the Swagger documentation will be updated.

To configure your API documentation, please refer to [swag's declarative comments format](https://github.com/swaggo/swag#declarative-comments-format) and [examples](https://github.com/swaggo/swag#examples).

To access the documentation, please visit http://localhost:8006/docs/index.html.
