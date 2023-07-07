// Package inserter provides some common inserter functions for
// mmdbwriter.Tree.
package inserter

import (
	"fmt"

	"github.com/maxmind/mmdbwriter/mmdbtype"
)

// Func is a function that returns the data type to be inserted into an
// mmdbwriter.Tree using some conflict resolution strategy.
type Func func(mmdbtype.DataType) (mmdbtype.DataType, error)

// FuncGenerator is a function that generates an Func given a
// value.
type FuncGenerator func(value mmdbtype.DataType) Func

// Remove any records for the network being inserted.
func Remove(_ mmdbtype.DataType) (mmdbtype.DataType, error) {
	return nil, nil
}

// ReplaceWith generates an inserter function that replaces the existing
// value with the new value.
func ReplaceWith(value mmdbtype.DataType) Func {
	return func(_ mmdbtype.DataType) (mmdbtype.DataType, error) {
		return value, nil
	}
}

// TopLevelMergeWith creates an inserter for Map values that will update an
// existing Map by adding the top-level keys and values from the new Map,
// replacing any existing values for the keys.
//
// Both the new and existing value must be a Map. An error will be returned
// otherwise.
func TopLevelMergeWith(newValue mmdbtype.DataType) Func {
	return func(existingValue mmdbtype.DataType) (mmdbtype.DataType, error) {
		newMap, ok := newValue.(mmdbtype.Map)
		if !ok {
			return nil, fmt.Errorf(
				"the new value is a %T, not a Map; TopLevelMergeWith only works if both values are Map values",
				newValue,
			)
		}

		if existingValue == nil {
			return newValue, nil
		}

		// A possible optimization would be to not bother copying
		// values that will be replaced.
		existingMap, ok := existingValue.(mmdbtype.Map)
		if !ok {
			return nil, fmt.Errorf(
				"the existing value is a %T, not a Map; TopLevelMergeWith only works if both values are Map values",
				existingValue,
			)
		}

		returnMap := existingMap.Copy().(mmdbtype.Map)

		for k, v := range newMap {
			returnMap[k] = v.Copy()
		}

		return returnMap, nil
	}
}

// DeepMergeWith creates an inserter that will recursively update an existing
// value. Map and Slice values will be merged recursively. Other values will
// be replaced by the new value.
func DeepMergeWith(newValue mmdbtype.DataType) Func {
	return func(existingValue mmdbtype.DataType) (mmdbtype.DataType, error) {
		return deepMerge(existingValue, newValue)
	}
}

func deepMerge(existingValue, newValue mmdbtype.DataType) (mmdbtype.DataType, error) {
	if existingValue == nil {
		return newValue, nil
	}
	if newValue == nil {
		return existingValue, nil
	}
	switch existingValue := existingValue.(type) {
	case mmdbtype.Map:
		newMap, ok := newValue.(mmdbtype.Map)
		if !ok {
			return newValue, nil
		}
		existingMap := existingValue.Copy().(mmdbtype.Map)
		for k, v := range newMap {
			nv, err := deepMerge(existingMap[k], v)
			if err != nil {
				return nil, err
			}
			existingMap[k] = nv
		}
		return existingMap, nil
	case mmdbtype.Slice:
		newSlice, ok := newValue.(mmdbtype.Slice)
		if !ok {
			return newValue, nil
		}
		length := len(existingValue)
		if len(newSlice) > length {
			length = len(newSlice)
		}

		rv := make(mmdbtype.Slice, length)
		for i := range rv {
			var ev, nv mmdbtype.DataType
			if i < len(existingValue) {
				ev = existingValue[i]
			}
			if i < len(newSlice) {
				nv = newSlice[i]
			}
			var err error
			rv[i], err = deepMerge(ev, nv)
			if err != nil {
				return nil, err
			}
		}
		return rv, nil
	default:
		return newValue, nil
	}
}
