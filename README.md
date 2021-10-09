# Mock Instagram Backend

This is an API created for a mock Instagram-like application for the internship recruitment task of Appointy in September, 2021 for Vellore Institute of Technology.
This API has been used using Golang and using `net/http` (along with other standard libraries) and no other external libraries or frameworks have been used as per the directions.

## Technologies Used

This backend is written in Golang and uses standard libraries for JSON encoding/decoding and `net/http` for routing and serving. The API is designed to communicate in JSON with its client.
MongoDB is used as the database as specified in the requirements of the task.

## Building the source
Ensure that Golang is installed and set up in your system.
First, clone this repository into a folder. Change directory to that folder.

```sh
git clone <this repository's URL.git>
cd AppointyTaskGo
```

Next, build the project. The API uses the Go driver for MongoDB. Building the project should automatically fetch the dependencies.
From within this directory (project root), run:

```sh
go build
```

## Running the API

To run the API, first set the environment variables `MONGODB_URI`, `MONGODB_DBNAME` and `APPYINSTA_PORT`.

### Linux

```sh
export MONGODB_URI=<your connection string>
export MONGODB_DBNAME=<your database name>
export APPYINST_PORT=<on which port to run the server>
```

### Windows (Powershell)
```ps
$Env:MONGODB_URI = "<your connection string>"
$Env:MONGODB_DBNAME = "<your database name>"
$Env:APPYINSTA_PORT = "<your connection string>"
```

You can also set these environment variables using other methods.

After this, run the executable created after building it.

### Linux
```sh
./appointy
```
You may have to change permissions for executing this.

### Windows (Powershell)
```ps
./appointy
```

The server will now run on the specified port (as specified in the `APPYINSTA_PORT` environment variable).

## API Specification

A simple overview of the API is as follows. The API has been designed and created as per the requirements specified in the task.


| Route | Method | Description | Request Body (sample) | Response Body |
| ------ | ------ | ------ | ------ | ------ |

