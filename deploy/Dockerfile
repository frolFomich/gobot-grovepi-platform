ARG BUILDER_IMAGE=golang:alpine
############################
# STEP 1 build executable binary
############################
FROM ${BUILDER_IMAGE} as builder

ENV APP_CMD=gobot-grovepi-platform

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache openssh git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/${APP_CMD}/

ENV GO111MODULE=on
# use modules
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build \
    -ldflags="-w -s -extldflags '-static'" -a \
    -v -o /go/bin/${APP_CMD} \
    ./cmd/grovepi/.

############################
# STEP 2 build a small image
############################
FROM alpine:latest

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable
COPY --from=builder /go/bin/gobot-grovepi-platform /go/bin/gobot-grovepi-platform

# Copy static content
WORKDIR /opt/config
COPY config .

# Use an unprivileged user.
USER appuser:998
WORKDIR /opt

# Run the binary.
ENTRYPOINT ["/go/bin/gobot-grovepi-platform", "-c", "/opt/config/app.yaml"]