# Simple Budget Application
## Description
This is a simple budget application that allows the user to add expenses and deposits to their budget then allocate the funds to track where the income is going. This application uses a Mongo database, . The application is downloadable and uses Mongo for the database. As of now, you will need to provide your own Mongo database to use this application. Add the mongo uri in a .env file in the same directory as the main.go file. The .env file should look like this:

```
MONGO_URI=<your mongo uri>
```

## Table of Contents
* [Installation](#installation)
* [Usage](#usage)
* [License](#license)
* [Questions](#questions)

## Installation
To install the application, clone the repository and run `make init` or `make build.win/lin` or `make run` to install the dependencies.

If you do not have make installed, you can run `go mod download` to install the dependencies. The application can be run with `go run main.go` or `go build main.go` and then `./main.exe` or `./main` depending on your operating system.

\* Note: The application will not run without a .env file with a mongo uri. For more info visit the official [MongoDB](https://www.mongodb.com/) website.

## Usage


