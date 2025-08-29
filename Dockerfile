FROM 416670754337.dkr.ecr.eu-west-2.amazonaws.com/ci-golang-build-1.24:latest
#FROM golang:1.24

WORKDIR /app

COPY . ./
COPY ecs-image-build/apispec ./apispec

CMD [ "/app/chs-delta-api" ]
