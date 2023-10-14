package template

var Response = `
package response

// {{.UpperCamlName}}Entity 实体
type {{.UpperCamlName}}Entity struct {
	// ID
	ID int32 {{.Backquote}}json:"id"{{.Backquote}}
}
`
