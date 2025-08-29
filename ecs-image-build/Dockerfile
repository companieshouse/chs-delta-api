FROM 416670754337.dkr.ecr.eu-west-2.amazonaws.com/ci-golang-build-1.24:latest

# Install curl for healthcheck
RUN apt-get update && apt-get install -y --no-install-recommends curl && rm -rf /var/lib/apt/lists/*

WORKDIR /opt

COPY apispec ./apispec
COPY /app .
COPY docker_start.sh .

RUN chmod +x docker_start.sh

CMD ["./docker_start.sh"]