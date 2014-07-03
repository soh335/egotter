package main

import (
	"reflect"
	"testing"
)

func CheckDuplicateMatch(t *testing.T, d *DuplicateMatch, data map[string]interface{}, expectMsg string, expectErr error, expectList []string) {
	msg, err := d.Match(data)
	if msg != expectMsg {
		t.Error("msg should be", expectMsg, ". but got ", msg)
	}
	if err != expectErr {
		t.Error("err should be", expectErr, ". but got", err)
	}
	if !reflect.DeepEqual(d.list, expectList) {
		t.Error("d.list should be", d.list, "but got", expectList)
	}
}

func TestDuplicateMatch(t *testing.T) {
	size := 2
	d := NewDuplicateMatch(size)

	CheckDuplicateMatch(
		t,
		d,
		map[string]interface{}{"id_str": "1"},
		"",
		nil,
		[]string{"1", ""},
	)

	CheckDuplicateMatch(
		t,
		d,
		map[string]interface{}{"id_str": "2"},
		"",
		nil,
		[]string{"2", "1"},
	)

	CheckDuplicateMatch(
		t,
		d,
		map[string]interface{}{"id_str": "3"},
		"",
		nil,
		[]string{"3", "2"},
	)

	CheckDuplicateMatch(
		t,
		d,
		map[string]interface{}{"id_str": "2"},
		"",
		NotChainError,
		[]string{"3", "2"},
	)

	CheckDuplicateMatch(
		t,
		d,
		map[string]interface{}{"hoge": 2},
		"",
		nil,
		[]string{"3", "2"},
	)
}
