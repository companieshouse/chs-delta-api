FROM 169942020521.dkr.ecr.eu-west-1.amazonaws.com/base/golang:1.16-alpine-builder AS golang-builder
FROM 169942020521.dkr.ecr.eu-west-1.amazonaws.com/base/golang:1.16-alpine-runtime AS golang-runtime
WORKDIR /app
COPY --from=golang-builder /build/apispec ./apispec
CMD ["-bind-addr=:5010"]
EXPOSE 5010

