FROM golang as build
WORKDIR /application
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w" -o app

FROM scratch
WORKDIR /application
LABEL "applicationName"="client"
COPY --from=build /application/app .
CMD ["./app"]