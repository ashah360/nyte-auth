FROM golang:1.17 as builder
LABEL maintainer="Arman Shah <ashah360@uw.edu>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ../main .

FROM alpine:latest
ARG DOPPLER_TOKEN
# Install the Doppler CLI
RUN (curl -Ls https://cli.doppler.com/install.sh || wget -qO- https://cli.doppler.com/install.sh) | sh
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
ENV PORT 3000
ENV DOPPLER_TOKEN ${DOPPLER_TOKEN}
EXPOSE 3000
CMD ["doppler", "run", "--", "./main"]