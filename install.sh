
GOBIN=$HOME/go/bin CGO_ENABLED=0 go install -ldflags "-X main.BUILD_INFO=pargocode:$(date '+%Y-%m-%d%n'):$(git rev-parse --short HEAD) -X main.AMBIENTE=PROD"