//go:build tools
// +build tools

package person

import (
	_ "github.com/axw/gocov/gocov"
	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "golang.org/x/tools/cmd/goimports"
	_ "gotest.tools/gotestsum"

	_ "github.com/google/wire/cmd/wire"
	_ "github.com/thoohv5/swagger-ci/cmd/swagger-ci"
	_ "golang.org/x/text/message"
	_ "golang.org/x/tools/cmd/stringer"
)
