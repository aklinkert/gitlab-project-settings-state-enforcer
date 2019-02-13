FROM scalify/glide:0.13.1 as builder
WORKDIR /go/src/github.com/Scalify/gitlab-project-settings-state-enforcer/

COPY glide.yaml glide.lock ./
RUN glide install --strip-vendor

COPY . ./
RUN CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o bin/gitlab-project-settings-state-enforcer .


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/Scalify/gitlab-project-settings-state-enforcer/bin/gitlab-project-settings-state-enforcer .
RUN chmod +x gitlab-project-settings-state-enforcer
ENTRYPOINT ["/root/gitlab-project-settings-state-enforcer"]
