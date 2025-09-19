package multi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type KeyValueInt64Flag struct {
	key   string
	value int64
}

func (e *KeyValueInt64Flag) Key() string {
	return e.key
}

func (e *KeyValueInt64Flag) Value() interface{} {
	return e.value
}

type KeyValueInt64 []*KeyValueInt64Flag

func (e *KeyValueInt64) String() string {
	return fmt.Sprintf("%v", *e)
}

func (e *KeyValueInt64) Set(value string) error {

	value = strings.Trim(value, " ")
	kv := strings.Split(value, SEP)

	if len(kv) != 2 {
		return errors.New("Invalid key=value argument")
	}

	v, err := strconv.ParseInt(kv[1], 10, 64)

	if err != nil {
		return err
	}

	a := KeyValueInt64Flag{
		key:   kv[0],
		value: v,
	}

	*e = append(*e, &a)
	return nil
}

func (e *KeyValueInt64) Get() interface{} {
	return *e
}
