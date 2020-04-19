package model

type Cpu struct {
	Idle float64
	Busy float64
}

type Disk struct {
	Free    uint64
	Total   uint64
	Used    uint64
	FsFile  string
	Percent float64
}

type LoadAvg struct {
	Name string
	Load float64
}

type Memory struct {
	Total uint64 //kb
	Used  uint64 //kb
	Free  uint64
}

type Port struct {
	Type   string
	Port   int64
	Listen bool
}
