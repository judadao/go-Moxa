package do

func NewMachine(main_type string, sub_type string, ip string, ch_type string, numOfChan int) *Machine {
	machine := &Machine{
		Main_type: main_type,
		Sub_type:  sub_type,
		IP:        ip,
		Ch_type:   ch_type,
		Channel:   make([]chan int, numOfChan),
		NumOfChan: numOfChan,
	}
	// 初始化每个通道
	for i := range machine.Channel {
		machine.Channel[i] = make(chan int, 1)
	}
	return machine
}

func NewMachine_4510(main_type string, sub_type string, ip string, slot_nick string, ch_type string, numOfChan int) *Machine {
	machine := &Machine{
		Main_type: main_type,
		Sub_type:  sub_type,
		Slot_nick: slot_nick,
		IP:        ip,
		Ch_type:   ch_type,
		Channel:   make([]chan int, numOfChan),
		NumOfChan: numOfChan,
	}
	// 初始化每个通道
	for i := range machine.Channel {
		machine.Channel[i] = make(chan int, 1)
	}
	return machine
}