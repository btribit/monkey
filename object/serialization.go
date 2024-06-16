// serialization.go
package object

import (
	"github.com/vmihailenco/msgpack"
)

// Serialize a FunctionCall to MsgPack format
func serializeFunctionCall(call FunctionCall) ([]byte, error) {
	return msgpack.Marshal(call)
}

// Deserialize MsgPack format to a FunctionCall
func deserializeFunctionCall(data []byte) (FunctionCall, error) {
	var call FunctionCall
	err := msgpack.Unmarshal(data, &call)
	return call, err
}

// Serialize a FunctionResponse to MsgPack format
func serializeFunctionResponse(resp FunctionResponse) ([]byte, error) {
	return msgpack.Marshal(resp)
}

// Deserialize MsgPack format to a FunctionResponse
func deserializeFunctionResponse(data []byte) (FunctionResponse, error) {
	var resp FunctionResponse
	err := msgpack.Unmarshal(data, &resp)
	return resp, err
}
