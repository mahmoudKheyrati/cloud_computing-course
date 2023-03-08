FROM golang as build
WORKDIR /application
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .

FROM scratch
LABEL "applicationName"="server"
WORKDIR /application
COPY --from=build /application/app .
EXPOSE 3000
CMD ["./app"]