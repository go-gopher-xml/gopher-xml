package xml

import (
	"encoding/xml"

	"github.com/yuin/gopher-lua"
)

var api = map[string]lua.LGFunction{
	"decode": apiDecode,
	"encode": apiEncode,
}

func apiDecode(L *lua.LState) int {
	str := L.CheckString(1)

	var value interface{}
	err := xml.Unmarshal([]byte(str), &value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(fromXML(L, value))
	return 1
}

func apiEncode(L *lua.LState) int {
	value := L.CheckAny(1)

	visited := make(map[*lua.LTable]bool)
	data, err := toXML(value, visited)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(data)))
	return 1
}
