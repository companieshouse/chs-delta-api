FROM golang:1.26

WORKDIR /app

COPY . ./
COPY ecs-image-build/apispec ./apispec

CMD [ "/app/chs-delta-api" ]
