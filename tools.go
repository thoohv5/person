//go:build tools
// +build tools

package person

import (
	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "golang.org/x/tools/cmd/goimports"

	_ "github.com/google/wire/cmd/wire"
	_ "golang.org/x/tools/cmd/stringer"
)
