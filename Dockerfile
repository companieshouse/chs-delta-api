FROM 169942020521.dkr.ecr.eu-west-1.amazonaws.com/base/golang:1.16-alpine-builder
FROM 169942020521.dkr.ecr.eu-west-1.amazonaws.com/base/golang:1.16-alpine-runtime
COPY --from=0 /build/apispec ./apispec
CMD ["-bind-addr=:5010"]
EXPOSE 5010

