FROM golang as build
WORKDIR /application
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o app -ldflags "-s -w"

FROM scratch
WORKDIR /application
LABEL "applicationName"="client"
COPY --from=build /application/app .
CMD ["./app"]