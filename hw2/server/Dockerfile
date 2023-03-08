FROM golang as build
WORKDIR /application
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o app -ldflags "-s -w"

FROM scratch
WORKDIR /application
LABEL "applicationName"="server"
COPY --from=build /application/app .
EXPOSE 3000
CMD ["./app"]