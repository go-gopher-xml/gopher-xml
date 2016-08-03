package xml

import (
	"github.com/yuin/gopher-lua"
)

// Preload adds xml to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local json = require("xml")
func Preload(L *lua.LState) {
	L.PreloadModule("xml", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}
