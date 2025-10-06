package multi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type KeyValueFloat64Flag struct {
	key   string
	value float64
}

func (e *KeyValueFloat64Flag) Key() string {
	return e.key
}

func (e *KeyValueFloat64Flag) Value() interface{} {
	return e.value
}

type KeyValueFloat64 []*KeyValueFloat64Flag

func (e *KeyValueFloat64) String() string {
	return fmt.Sprintf("%v", *e)
}

func (e *KeyValueFloat64) Set(value string) error {

	value = strings.Trim(value, " ")
	kv := strings.Split(value, SEP)

	if len(kv) != 2 {
		return errors.New("Invalid key=value argument")
	}

	v, err := strconv.ParseFloat(kv[1], 64)

	if err != nil {
		return err
	}

	a := KeyValueFloat64Flag{
		key:   kv[0],
		value: v,
	}

	*e = append(*e, &a)
	return nil
}

func (e *KeyValueFloat64) Get() interface{} {
	return *e
}
