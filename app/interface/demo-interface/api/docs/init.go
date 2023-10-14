// Package docs DOCS
package docs

import "github.com/swaggo/swag"

// SetBasePath 设置基础路径
func SetBasePath(name, bp string) {
	spec, ok := swag.GetSwagger(name).(*swag.Spec)
	if ok {
		spec.BasePath = bp
	}
}

// SetVersion 设置版本
func SetVersion(name, version string) {
	spec, ok := swag.GetSwagger(name).(*swag.Spec)
	if ok {
		spec.Version = version
	}
}
