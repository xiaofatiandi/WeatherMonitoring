# Weather Monitoring Service

# Description
This service provides these functionalities:
  - Enroll a device
  - Submit latest temperature readings from the enrolled device
  - Fetch aggregated temperature values.


# Endpoints

- `POST /devices/{id}` : enroll a device
- `GET /devices"` : list all devices
- `POST /device/enable/{id}` : enable a device 
- `POST /device/disable/{id}` : disable a device
- `POST /temperature`: post a temperature with device_id and temperature value
- `GET /temperature/aggregated`: get aggregated temperature of a device for all devices


# How to run 
You may either run the service on local or on docker.  
## Run on local
Requirement: need to install go sdk 

To start the service, in a terminal, run :
```
cd src
go run .
```

## Run on docker
I have created the `/deploymen/Dockerfile` for building the docker image. To build the docker image:
```
docker build -t weather-monitoring -f deployment/Dockerfile .
```
To run the docker image:
```
docker run -p 8000:8000 weather-monitoring                   
```

## To query the endpoints of the service
You may choose various clients. Below shows the examples to use curl commands:

- Enroll a device:
```
curl -X POST http://localhost:8000/devices/804baf2e-da41-4e37-b81d-f16763f0f932        
```

- List all devices
```
curl -X GET http://localhost:8000/devices                                            
```

- Post a temperature from a device
```
curl -X POST http://localhost:8000/temperature \
     -d '{"device_id": "804baf2e-da41-4e37-b81d-f16763f0f932", "temperature": 25.4}'
```

- Get aggregated temperature of all devices
```
curl -X GET http://localhost:8000/temperature/aggregated
```



## Run test
```
go test -v ./...
```


# Deploy 
I have created the Dockerfile for building docker image. Build the image and deploy to the desired cloud environments.


# Notes to the reviewers

## Folder structure
```
/src
    /handler //define web request handlers
        handler.go
        handler_test.go
    /logger //define logger to log info, error etc in desired format.
        logger.go
    /storage //define storage for storing data and retrieving data
        storage.go
        storage_test.go
        models.go //define data structures
    go.mod
    go.sum

/deployment //deployment files
    Dockerfile
        
```


## For design, I divide the service implementation into three main components:

### 1. Server runner (`main.go`)
  For server runner, I use `github.com/gorilla/mux`, there are some other libraries that can be used such as `net/http` etc. I choose mux because it provides easy extraction of variables from the URL, and easy routes. 

### 2. Web request handler (`/handler`)
  
  
  `handler` class  defines the functions to handle all web requests which will call `Storage` interface to store and retrieve data

 ### 3. Storage ( `/storage`)
  
  `Storage` interface decouples the code from specific implementation. This makes the code more flexible and easier to maintain. For this example, I implement a simple `InMemoryStorage` in this project for now as I don't have a local database and want to make it simple. But if in the future, we need to switch to use database, we only need to create a new implementation of the `Storage` interface with database, and no need to change any other code. 


I also implement a `logger` class in `/logger` to support better logging. 


## Future improvements:
One  thing that I didn't do here is that currently I hard coded the port and server ip using localhost:8000, these can be moved to a configuration file, and main.go should read those values from configuration file, and docker deployment should also use the same values 
