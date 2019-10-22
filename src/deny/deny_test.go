package deny

//import (
//	"log"
//	"testing"
//	"time"
//)
//
//func TestInjection(t *testing.T) {
//	type obj struct {
//		Str      string
//		Age      int32
//		Address  []string
//		CreateAt time.Time
//		Fun      float64
//		Ddk      []byte
//	}
//	nobj := &obj{
//		Str:      "jstefdfdf",
//		Age:      2,
//		Address:  []string{"kjdfkjdkf", "kfldkfjdf$"},
//		CreateAt: time.Now(),
//		Fun:      4545.5,
//		Ddk:      []byte{'d','f'},
//	}
//
//	res := Injection(nobj)
//	log.Println(*res)
//	if *res != false {
//		t.Error("no pass")
//	}
//}
