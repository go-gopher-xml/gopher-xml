package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"

	"github.com/yuin/gopher-lua"
)

var (
	errFunction = errors.New("cannot encode function to XML")
	errChannel  = errors.New("cannot encode channel to XML")
	errState    = errors.New("cannot encode state to XML")
	errUserData = errors.New("cannot encode userdata to XML")
	errNested   = errors.New("cannot encode recursively nested tables to XML")
)

type xmlValue struct {
	lua.LValue
	visited map[*lua.LTable]bool
}

func (j xmlValue) MarshalXML() ([]byte, error) {
	return toXML(j.LValue, j.visited)
}

func toXML(value lua.LValue, visited map[*lua.LTable]bool) (data []byte, err error) {
	switch converted := value.(type) {
	case lua.LBool:
		data, err = xml.Marshal(converted)
	case lua.LChannel:
		err = errChannel
	case lua.LNumber:
		data, err = xml.Marshal(converted)
	case *lua.LFunction:
		err = errFunction
	case *lua.LNilType:
		data, err = xml.Marshal(converted)
	case *lua.LState:
		err = errState
	case lua.LString:
		data, err = xml.Marshal(converted)
	case *lua.LTable:
		var arr []xmlValue
		var obj map[string]xmlValue

		if visited[converted] {
			panic(errNested)
		}
		visited[converted] = true

		converted.ForEach(func(k lua.LValue, v lua.LValue) {
			if k, ok := k.(lua.LString); ok {
				fmt.Println("k:" + string(k))
			}
			if v, ok := v.(lua.LNumber); ok {
				fmt.Println("v:" + strconv.Itoa(int(v)))
			}
			// fmt.Println("v" + v)
			i, numberKey := k.(lua.LNumber)
			if numberKey && obj == nil {
				index := int(i) - 1
				if index != len(arr) {
					// map out of order; convert to map
					obj = make(map[string]xmlValue)
					for i, value := range arr {
						obj[strconv.Itoa(i+1)] = value
					}
					obj[strconv.Itoa(index+1)] = xmlValue{v, visited}
					return
				}
				arr = append(arr, xmlValue{v, visited})
				return
			}
			if obj == nil {
				obj = make(map[string]xmlValue)
				for i, value := range arr {
					obj[strconv.Itoa(i+1)] = value
				}
			}
			obj[k.String()] = xmlValue{v, visited}
		})
		if obj != nil {
			fmt.Printf("obj:%v", obj)
			data, err = xml.Marshal(obj)
		} else {
			data, err = xml.Marshal(arr)
		}
		fmt.Println("data" + string(data))
	case *lua.LUserData:
		// TODO: call metatable __tostring?
		err = errUserData
	}
	return
}

func fromXML(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(fromXML(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), fromXML(L, item))
		}
		return tbl
	}
	return lua.LNil
}
