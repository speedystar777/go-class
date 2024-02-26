# STAGE 1:
# first stage is actual build process, starts by using a golang container
FROM golang:1.21-alpine AS builder 

# pulling in certificates
RUN /sbin/apk update && \
    /sbin/apk --no-cache add ca-certificates git tzdata && \
    /usr/sbin/update-ca-certificates

# adding a user specific to our app and picking a work directory to work in
RUN adduser -D -g '' sort
WORKDIR /home/sort

# copying what we need to build program into our container
COPY go.mod /home/sort
COPY go.sum /home/sort
COPY cmd    /home/sort/cmd
COPY *.go   /home/sort

# define version argument
ARG VERSION

# command to build program as binary -> version argument is baked into executable
RUN CGO_ENABLED=0 go build -a -tags netgo,osusergo -ldflags "-extldflags '-static' -s -w" -ldflags "-X main.version=$VERSION" -o sort ./cmd/sort

# STAGE 2:
# busybox is package that provides basic unix tools in small, statically linked package
# better than having a lot of os executables
FROM busybox:musl

# copying certificates, timezone info, password for user, and sort binary 
# when you copy the sort binary you are copying the whole program that you wrote
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /home/sort/sort /home/sort

# the user that we added
USER sort
# the home directory
WORKDIR /home
# the port to expose
EXPOSE 8081

# entrypoint to start binary
ENTRYPOINT ["/home/sort"]