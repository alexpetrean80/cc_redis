package resp_test

import (
	"fmt"
	"testing"

	"github.com/alexpetrean80/cc_redis/resp"
)

// test the following scenarios:
//  "$-1\r\n"
//  "*1\r\n$4\r\nping\r\n”
//  "*2\r\n$4\r\necho\r\n$5\r\nhello world\r\n”
//  "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n”
//  "+OK\r\n"
//  "-Error message\r\n"
//  "$0\r\n\r\n"
//  "+hello world\r\n”

func TestNull(t *testing.T) {
	resp := resp.New([]byte("$-1\r\n"))
	v, err := resp.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if v != nil {
		t.Fatal("Expected nil")
	}
}

func TestSimpleString(t *testing.T) {
	resp := resp.New([]byte("+OK\r\n"))
	v, err := resp.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if v != "OK" {
		t.Fatal("Expected OK")
	}
}

func TestErrorMsg(t *testing.T) {
	resp := resp.New([]byte("-Error message\r\n"))
	v, err := resp.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if v.(error).Error() != "Error message" {
		t.Fatal("Expected Error message")
	}
}

func TestEmptyStr(t *testing.T) {
	resp := resp.New([]byte("$0\r\n\r\n"))
	v, err := resp.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if v != "" {
		t.Fatal("Expected empty string")
	}
}

func TestStr(t *testing.T) {
	resp := resp.New([]byte("$11\r\nhello world\r\n"))
	v, err := resp.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if v != "hello world" {
		t.Fatal("Expected hello world")
	}
}

func TestArrSimple(t *testing.T) {
	resp := resp.New([]byte("*1\r\n+ping\r\n"))
	v, err := resp.Parse()
	fmt.Println("v: ", v)
	if err != nil {
		t.Fatal(err)
	}
	if len(v.([]interface{})) == 0 {
		t.Fatal("Empty array")
	}
	if v.([]interface{})[0].(string) != "ping" {
		t.Fatal("Expected ping")
	}
}

func TestArr1(t *testing.T) {
	resp := resp.New([]byte("*1\r\n$4\r\nping\r\n"))
	v, err := resp.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(v.([]interface{})) == 0 {
		t.Fatal("Empty array")
	}

	if v.([]interface{})[0] != "ping" {
		t.Fatal("Expected ping")
	}
}

func TestArr2(t *testing.T) {
	resp := resp.New([]byte("*2\r\n$4\r\necho\r\n$5\r\nhello world\r\n"))
	v, err := resp.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(v.([]interface{})) == 0 {
		t.Fatal("Empty array")
	}

	expected := []string{"echo", "hello world"}
	if len(v.([]interface{})) != len(expected) {
		t.Fatal("wrong length")
	}
	for i := 0; i < len(expected); i++ {
		if v.([]interface{})[i] != expected[i] {
			t.Fatalf("%s not %s", v.([]interface{})[i], expected[i])
		}
	}
}

func TestArr3(t *testing.T) {
	resp := resp.New([]byte("*2\r\n$3\r\nget\r\n$3\r\nkey\r\n"))
	v, err := resp.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(v.([]interface{})) == 0 {
		t.Fatal("Empty array")
	}

	if v != "key" {
		t.Fatal("Expected key")
	}
}
