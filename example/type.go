package main

type Apple struct {
}

func (a Apple) Color() string {
	return "Red"
}

func (a Apple) Rotten() bool {
	return false
}
