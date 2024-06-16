// protocol.go
package object

type FunctionCall struct {
	Type string        `msgpack:"type"`
	Name string        `msgpack:"name"`
	Args []interface{} `msgpack:"args"`
}

type FunctionResponse struct {
	Type   string      `msgpack:"type"`
	Result interface{} `msgpack:"result"`
	Error  *string     `msgpack:"error"`
}
