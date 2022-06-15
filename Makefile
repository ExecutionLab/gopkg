.PHONY: lint gci

lint:
	@hash golangci-lint 2>/dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.42.1
	@golangci-lint run -v

gci:
	@GO111MODULE=off go get github.com/daixiang0/gci
	gci -w -local gitlab.com/husol/pkg .
