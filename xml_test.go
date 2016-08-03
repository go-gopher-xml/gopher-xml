package xml

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

func TestSimple(t *testing.T) {
	const str = `
	local xml = require("xml")
	assert(type(xml) == "table")
	assert(type(xml.decode) == "function")
	assert(type(xml.encode) == "function")
	local obj = {["Ac"]=1,["B"]=2,["C"]=3}
    print(obj["a"])
	local xmlStr = xml.encode(obj)
    print(xmlStr)
    print(xml.encode(true))
	assert(xml.encode(true) == "true")
	assert(xml.encode(1) == "1")
	assert(xml.encode(-10) == "-10")
	assert(xml.encode(nil) == "{}")

	local obj = {["a"]=1,["b"]=2,["c"]=3}
	local xmlStr = xml.encode(obj)
    print(xmlStr)
	local xmlObj = xml.decode(xmlStr)
	-- for i = 1, #obj do
	assert(obj["a"] == xmlObj["a"])
	-- end

	local obj = {name="Tim",number=12345}
	local xmlStr = xml.encode(obj)
	local xmlObj = xml.decode(xmlStr)
	assert(obj.name == xmlObj.name)
	assert(obj.number == xmlObj.number)

	local obj = {"a","b",what="c",[5]="asd"}
	local xmlStr = xml.encode(obj)
	local xmlObj = xml.decode(xmlStr)
	assert(obj[1] == xmlObj["1"])
	assert(obj[2] == xmlObj["2"])
	assert(obj.what == xmlObj["what"])
	assert(obj[5] == xmlObj["5"])

	assert(xml.decode("null") == nil)

	assert(xml.decode(xml.encode({person={name = "tim",}})).person.name == "tim")

	local obj = {
		abc = 123,
		def = nil,
	}
	local obj2 = {
		obj = obj,
	}
	obj.obj2 = obj2
	assert(xml.encode(obj) == nil)
	`
	s := lua.NewState()
	Preload(s)
	if err := s.DoString(str); err != nil {
		t.Error(err)
	}
}

func TestCustomRequire(t *testing.T) {
	const str = `
	local j = require("xml")
	assert(type(j) == "table")
	assert(type(j.decode) == "function")
	assert(type(j.encode) == "function")
	`
	s := lua.NewState()
	s.PreloadModule("xml", Loader)
	if err := s.DoString(str); err != nil {
		t.Error(err)
	}
}
