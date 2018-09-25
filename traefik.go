package main

type TreafikConf struct {
	EtcdPrefix string
}

type TraefikBackend struct {
	Name   string
	Node   string
	URL    string
	Weight uint
}
