package main

type TalkType struct {
	Id   uint
	Name string
}

type Talk struct {
	Id          uint32
	Presenter   string
	TypeId      TalkType
	Name        string
	Description string
	Type        TalkType
}
