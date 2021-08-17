# Debugging the chs-delta-api

## Overview

When using the chs-delta-api within the docker-chs-development environment, we can make use of the Golang debug builder 
and runtime images to create a debug image for this service. This debug image can be used in docker-chs-development to 
debug the service during runtime using the Goland IDE debugging tools.

## 1. Enabling debugging in chs-delta-api

### 1.1 Creating the temporary debug image
Open the Dockerfile found at the root of the chs-delta-api project and replace it with the following
```Dockerfile
FROM 169942020521.dkr.ecr.eu-west-2.amazonaws.com/base/golang:1.16-alpine-debug-builder as debug-builder
FROM 169942020521.dkr.ecr.eu-west-2.amazonaws.com/base/golang:1.16-alpine-debug-runtime as debug-runtime

WORKDIR /app
COPY --from=debug-builder /build/apispec ./apispec

CMD ["./dlv", "exec", "--listen=:5010", "--headless=true", "--api-version=2", "./app"]

EXPOSE 5010
```

The above Dockerfile uses the previously mentioned debug builder and runtime images to create a base image for our file. 
We then copy over the api spec from the temporary build folder into our final image folder space as this spec is required for 
the application to work correctly.

Finally, we set the CMD to the included dlv tool, providing all the necessary arguments to allow dlv to start our service.

### 1.2 Mapping the docker-compose ports on our local machine
Open the docker-chs-development repository in an IDE and edit the `chs-delta-api.docker-compose.yaml` file and add the 
following properties:
```yaml
services:
  chs-delta-api:
    ...
    environment:
      ...
      - BIND_ADDR=:5010
    ports:
      - "5010:5010"
```

This allows us to use our local machine's port to connect to the matched docker container port.

Make sure the docker port you select matches the port you previously exposed in the Dockerfile.

## 2. Using the debugger on the chs-delta-api

### 2.1 Starting the debugger in Goland IDE
Now we've started our application using the delve debugging tool we need to attach our debugger to delve. To do this 
we need to open our project in Goland IDE and visit the `Run -> Edit Configurations` window. Press the `+` to add a new 
configuration and select `Go Remote`. Set the following values:
- `Name`: Go debug
- `Host`: api.chs.local
- `Port`: 5010

Press `Apply` and finally in the top right-hand corner of Goland IDE select your new connection and press `Debug`.

### 2.2 Communicating with the service while in debug mode
As this is an internal API, all traffic routes through ERIC. When using docker-chs-development we will want to use 
the following URL to communicate with the chs-delta-api service: `http://api.chs.local:4001/delta/officers`. Now all you 
need to do is add your breakpoints and debug as you would any other service.

## Final notes
This image shouldn't be committed to GitHub or AWS ECR.

Documentation on delve can be found here: https://github.com/go-delve/delve
