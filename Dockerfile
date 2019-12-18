FROM golang:1.13.4-alpine AS build-env
ADD . /src
RUN cd /src && go build -o goapp

FROM alpine
ARG PORT=$PORT
WORKDIR /app
COPY --from=build-env /src/goapp /app/
CMD /app/goapp
