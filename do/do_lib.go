package do

type DoRequest struct {
	APIkey string
	IP     string
	Ch     int
}

type do_allApiFunc func(string, *Machine, int, string) int

var doApiMap = map[string]do_allApiFunc{
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

// key must same with doApiMap
var restUri = map[string]string{
	"DO_WHOLE":     "/api/slot/0/io/do",
	"DO_GET_VALUE": "/api/slot/0/io/do",
	"DO_PUT_VALUE": "/api/slot/0/io/do",
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
	IP        string
	Ch_type   string
	Channel   []chan int
	NumOfChan int
}
