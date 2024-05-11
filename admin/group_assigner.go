package admin

type groupAssigner struct {
	pendingTrips chan NewTripInfo
	newGroupIds  chan int64
	resultChan   chan NewGroupId
}

// AssignGroup implements GroupAssigner.
func (g *groupAssigner) AssignGroup(trip NewTripInfo) {
	go func() {
		g.pendingTrips <- trip
	}()
}

func (g *groupAssigner) OnNewGroup(groupId int64) {
	go func() {
		g.newGroupIds <- groupId
	}()
}

func (g *groupAssigner) ResultChan() chan NewGroupId {
	return g.resultChan
}

type GroupAssigner interface {
	AssignGroup(trip NewTripInfo)
	OnNewGroup(link int64)
	ResultChan() chan NewGroupId
}

func NewGroupAssigner() GroupAssigner {
	g := &groupAssigner{
		pendingTrips: make(chan NewTripInfo),
		newGroupIds:  make(chan int64),
		resultChan:   make(chan NewGroupId),
	}
	go func() {
		for {
			trip := <-g.pendingTrips
			grpId := <-g.newGroupIds
			g.resultChan <- NewGroupId{GroupId: grpId, NewTripInfo: trip}
		}
	}()
	return g

}
