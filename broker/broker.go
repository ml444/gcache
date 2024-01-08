package broker

/*
broker manage groups -> group manage shardMap -> shard manage entries
*/

type Broker struct {
	Address  string
	GroupMap map[string]*Group
}

func NewBroker(address string) *Broker {
	return &Broker{
		Address:  address,
		GroupMap: make(map[string]*Group),
	}
}

func (c *Broker) AddGroup(group *Group) {
	c.GroupMap[group.Name] = group
}
