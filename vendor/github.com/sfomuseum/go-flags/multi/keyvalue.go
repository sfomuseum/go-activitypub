package multi

import (
	"errors"
	"fmt"
	"strings"
)

const SEP string = "="

type KeyValueFlag interface {
	Key() string
	Value() interface{}
}

type KeyValueStringFlag struct {
	KeyValueFlag
	key   string
	value string
}

func (e *KeyValueStringFlag) Key() string {
	return e.key
}

func (e *KeyValueStringFlag) Value() interface{} {
	return e.value
}

type KeyValueCSVString []*KeyValueStringFlag

func (e *KeyValueCSVString) String() string {

	parts := make([]string, len(*e))

	for idx, k := range *e {
		parts[idx] = fmt.Sprintf("%s=%s", k.Key(), k.Value().(string))
	}

	return strings.Join(parts, ",")
}

func (e *KeyValueCSVString) Set(value string) error {

	for _, v := range strings.Split(value, ",") {

		value = strings.Trim(v, " ")
		kv := strings.Split(v, SEP)

		if len(kv) != 2 {
			return errors.New("Invalid key=value argument")
		}

		a := KeyValueStringFlag{
			key:   kv[0],
			value: kv[1],
		}

		*e = append(*e, &a)
	}

	return nil
}

type KeyValueString []*KeyValueStringFlag

func (e *KeyValueString) String() string {
	return fmt.Sprintf("%v", *e)
}

func (e *KeyValueString) Set(value string) error {

	value = strings.Trim(value, " ")
	kv := strings.Split(value, SEP)

	if len(kv) != 2 {
		return errors.New("Invalid key=value argument")
	}

	a := KeyValueStringFlag{
		key:   kv[0],
		value: kv[1],
	}

	*e = append(*e, &a)
	return nil
}

func (e *KeyValueString) Get() interface{} {
	return *e
}
