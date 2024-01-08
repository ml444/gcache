package node

import (
	"io/ioutil"

	"github.com/hashicorp/memberlist"
)

func New(addr, cluster string) (Node, error) {
	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = ioutil.Discard
	l, e := memberlist.Create(conf)
	if e != nil {
		return nil, e
	}
	if cluster == "" {
		cluster = addr
	}
	clu := []string{cluster}
	_, e = l.Join(clu)
	if e != nil {
		return nil, e
	}
	//circle := consistent.New()
	//circle.NumberOfReplicas = 256
	//go func() {
	//	for {
	//		m := l.Members()
	//		nodes := make([]string, len(m))
	//		for i, n := range m {
	//			nodes[i] = n.Name
	//		}
	//		circle.Set(nodes)
	//		time.Sleep(time.Second)
	//	}
	//}()
	return &node{circle, addr}, nil
}
