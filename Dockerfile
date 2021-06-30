FROM 169942020521.dkr.ecr.eu-west-2.amazonaws.com/base/golang:1.16-alpine-builder
FROM 169942020521.dkr.ecr.eu-west-2.amazonaws.com/base/golang:1.16-alpine-runtime
CMD ["-bind-addr=:5010"]
EXPOSE 5010

