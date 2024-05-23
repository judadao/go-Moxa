package do

type DoRequest struct {
	APIkey string
	IP     string
	Ch     int
}

type MachineType1200Func func(string, *Machine, int, string) int
type MachineType4510Func func(string, *Machine, string, string) int

type MachineTypeFunc interface {
	Call(string, *Machine, interface{}, string) int
}

type MachineType1200 struct{}
type MachineType4510 struct{}

var doApiMap = map[string]MachineType1200Func{
	"DO_WHOLE":          do_update_value,
	"DO_CHECK":          do_check_value,
	"DO_STATUS":         do_get_status,
	"DO_PAULSESTATUS":   do_get_paulse_status,
	"DO_PAULSECOUNT":    do_get_paulse_count,
	"DO_UPDATE_VALUE":   do_update_value,
	"DO_GET_VALUE":      do_get_value,
	"DO_PUT_VALUE":      do_put_value,
	"TEST_DO_GET_VALUE": test_do_get_value,
	"TEST_DO_PUT_VALUE": test_do_put_value,
}

var do4510_ApiMap = map[string]MachineType4510Func{
	"DO_GET_VALUE": do4510_get_value,
	// "DO_CHECK":     do_check_value,
	// "DO_GET_VALUE": do_get_value,
	// "DO_PUT_VALUE": do_put_value,
}

// key must same with doApiMap
var restUri = map[string]string{
	"DO_WHOLE":     "/api/slot/0/io/do",
	"DO_GET_VALUE": "/api/slot/0/io/do",
	"DO_PUT_VALUE": "/api/slot/0/io/do",
}

var Subtype_map = map[string]int{
	"e1211": 16,
	"e1212": 8,
	"e1213": 8,
	"e1214": 8,
	"e1242": 4,
}

var restParam = map[string]string{
	"DO_WHOLE":        "",
	"DO_STATUS":       "/doStatus",
	"DO_PULSE_COUNT":  "/doPulseCount",
	"DO_PULSE_STATUS": "/doPulseStatus",
}

type Machine struct {
	Main_type string
	Sub_type  string
	Slot_nick string
	IP        string
	Ch_type   string
	Channel   []chan int
	NumOfChan int
}
