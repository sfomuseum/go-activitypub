package multi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type KeyValueBoolFlag struct {
	key   string
	value bool
}

func (e *KeyValueBoolFlag) Key() string {
	return e.key
}

func (e *KeyValueBoolFlag) Value() interface{} {
	return e.value
}

type KeyValueBool []*KeyValueBoolFlag

func (e *KeyValueBool) String() string {
	return fmt.Sprintf("%v", *e)
}

func (e *KeyValueBool) Set(value string) error {

	value = strings.Trim(value, " ")
	kv := strings.Split(value, SEP)

	if len(kv) != 2 {
		return errors.New("Invalid key=value argument")
	}

	v, err := strconv.ParseBool(kv[1])

	if err != nil {
		return err
	}

	a := KeyValueBoolFlag{
		key:   kv[0],
		value: v,
	}

	*e = append(*e, &a)
	return nil
}

func (e *KeyValueBool) Get() interface{} {
	return *e
}
