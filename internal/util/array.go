package util

import (
	"net"
	"reflect"
	"strconv"
)

func InArrStr(slice []string, val string) (flag bool) {
	for _, item := range slice {
		if item == val {
			flag = true
			break
		}
	}
	return
}

func InArrInt64(slice []int64, val int64) (flag bool) {
	for _, item := range slice {
		if item == val {
			flag = true
			break
		}
	}
	return
}

func InArrInt32(slice []int32, val int32) (flag bool) {
	for _, item := range slice {
		if item == val {
			flag = true
			break
		}
	}
	return
}

func InArrIP(slice []net.Addr, val net.Addr) (flag bool) {
	for _, item := range slice {
		if item == val {
			flag = true
			break
		}
	}
	return
}

func StrArr2IntArr(strArr []string) (intArr []int64, err error) {
	for _, item := range strArr {
		if intVal, err := strconv.ParseInt(item, 10, 32); err != nil {
			return nil, err
		} else {
			intArr = append(intArr, intVal)
		}
	}
	return
}

func IntArr2StrArr(intArr []int32) (strArr []string) {
	for _, item := range intArr {
		strArr = append(strArr, strconv.FormatInt(int64(item), 10))
	}
	return
}

func UniqInt32Arr(intArr []int32) (uniqArr []int32) {
	uniqArr = []int32{}
	intMap := make(map[int32]bool)
	for _, intVal := range intArr {
		intMap[intVal] = true
	}

	for key, _ := range intMap {
		uniqArr = append(uniqArr, key)
	}

	return uniqArr
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func UniqStrArr(strArr []string) (uniqArr []string) {
	uniqArr = []string{}
	intMap := make(map[string]bool)
	for _, intVal := range strArr {
		intMap[intVal] = true
	}

	for key, _ := range intMap {
		uniqArr = append(uniqArr, key)
	}

	return uniqArr
}
