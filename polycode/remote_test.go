package polycode

//import (
//	"encoding/json"
//	"testing"
//)
//
//type Foo struct {
//	ID  string `json:"id"`
//	Msg string `json:"msg"`
//}
//
//type TaskOutput2 struct {
//	IsAsync bool   `json:"isAsync"`
//	IsNull  bool   `json:"isNull"`
//	Output  any    `json:"output"`
//	Error   *Error `json:"error"`
//}
//
//func TestFuture(t *testing.T) {
//	s := "{\"id\":\"\\\"test-1726297871495888\\\"\",\"msg\":\"\"}"
//	j := FutureImpl{data: s}
//
//	var ret string
//	err := j.Get(&ret)
//	if err != nil {
//		t.Errorf("Error: %v", err)
//	}
//	if ret != "test-1726297871495888" {
//		t.Errorf("Error: %v", ret)
//	}
//}
//
//func TestTaskOutput(t *testing.T) {
//	f := Foo{
//		ID:  "test-1726297871495888",
//		Msg: "",
//	}
//	output := TaskOutput2{
//		IsAsync: false,
//		IsNull:  false,
//		Output:  f,
//		Error:   nil,
//	}
//
//	b, err := json.Marshal(output)
//	if err != nil {
//		t.Errorf("Error: %v", err)
//	}
//
//	ret := TaskOutput2{}
//	err = json.Unmarshal(b, &ret)
//	println(ret.Output)
//
//	future := FutureFrom(ret.Output)
//
//	var ret2 string
//	err = future.Get(&ret2)
//	if err != nil {
//		t.Errorf("Error: %v", err)
//	}
//	if ret2 != "test-1726297871495888" {
//		t.Errorf("Error: %v", ret2)
//	}
//
//}
