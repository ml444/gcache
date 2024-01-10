package node

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"
)

//func New(addr, cluster string) (Node, error) {
//	conf := memberlist.DefaultLANConfig()
//	conf.Name = addr
//	conf.BindAddr = addr
//	conf.LogOutput = io.Discard
//	l, e := memberlist.Create(conf)
//	if e != nil {
//		return nil, e
//	}
//	if cluster == "" {
//		cluster = addr
//	}
//	clu := []string{cluster}
//	_, e = l.Join(clu)
//	if e != nil {
//		return nil, e
//	}
//	//circle := consistent.New()
//	//circle.NumberOfReplicas = 256
//	//go func() {
//	//	for {
//	//		m := l.Members()
//	//		nodes := make([]string, len(m))
//	//		for i, n := range m {
//	//			nodes[i] = n.Name
//	//		}
//	//		circle.Set(nodes)
//	//		time.Sleep(time.Second)
//	//	}
//	//}()
//	return &node{circle, addr}, nil
//}

var (
	mtx        sync.RWMutex
	items      = map[string]string{}
	broadcasts *memberlist.TransmitLimitedQueue

	join = flag.String("join", "", "comma seperated list of members")
	port = flag.Int("port", 4001, "http port")
)

func init() {
	flag.Parse()
}

type update struct {
	Action string
	Data   map[string]string
}

type broadcast struct {
	msg    []byte
	notify chan<- struct{}
}

func (b *broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (b *broadcast) Message() []byte {
	return b.msg
}

func (b *broadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}
func addHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	key := r.Form.Get("key")
	val := r.Form.Get("val")
	mtx.Lock()
	items[key] = val
	mtx.Unlock()
	b, err := json.Marshal([]*update{
		{
			Action: "add",
			Data: map[string]string{
				key: val,
			},
		},
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// 广播数据
	broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})
}

func delHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	key := r.Form.Get("key")
	mtx.Lock()
	delete(items, key)
	mtx.Unlock()

	b, err := json.Marshal([]*update{
		{
			Action: "del",
			Data: map[string]string{
				key: "",
			},
		},
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	key := r.Form.Get("key")
	mtx.RLock()
	val := items[key]
	mtx.RUnlock()
	if _, err := w.Write([]byte(val)); err != nil {
		log.Printf("fail to write response, err: %s.\n", err)
	}
}

func start(members []string) error {
	hostname, _ := os.Hostname()
	c := memberlist.DefaultLocalConfig()
	c.Delegate = &delegate{}
	c.BindPort = 0
	uid, _ := uuid.NewUUID()
	c.Name = hostname + "-" + uid.String()
	// 创建 Gossip 网络
	m, err := memberlist.Create(c)
	if err != nil {
		return err
	}
	// 第一个节点没有 member，从第二个开始有 member
	if len(members) > 0 {
		_, err := m.Join(members)
		if err != nil {
			return err
		}
	}
	broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return m.NumMembers()
		},
		RetransmitMult: 3,
	}
	node := m.LocalNode()
	log.Printf("Local member %s:%d\n", node.Addr, node.Port)
	return nil
}

func main() {
	mlist := strings.Split(*join, ",")
	if err := start(mlist); err != nil {
		panic(err)
	}
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/del", delHandler)
	http.HandleFunc("/get", getHandler)
	fmt.Printf("Listening on :%d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(err)
	}
}
