package helpers

import "reflect"

//InArray -> search element in array
func InArray(val interface{}, array interface{}) (index int) {
	values := reflect.ValueOf(array)
	if reflect.TypeOf(array).Kind() == reflect.Slice || values.Len() > 0 {
		for i := 0; i < values.Len(); i++ {
			if reflect.DeepEqual(val, values.Index(i).Interface()) {
				return i
			}
		}
	}
	return -1
}
