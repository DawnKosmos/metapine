package live

import "fmt"

type UpdateMode int

const ONCLOSE UpdateMode = 0

var ONTICK = -1

/*
The UpdateGroup manages that  calculations of an Indicator get Updated in the right order.
*/
type UpdateGroup struct {
	Name string
	ug   []Updater
	re   int64
	tick chan bool
	exit bool
}

func NewExpertProgramm(Name string, resolution int64) *UpdateGroup {
	return &UpdateGroup{
		Name: "",
		ug:   nil,
		re:   0,
	}
}

func (ug *UpdateGroup) AppendUpdater(u Updater) {
	ug.ug = append(ug.ug, u)
}

func (ug *UpdateGroup) Resolution() int64 {
	return ug.re
}

func (ug *UpdateGroup) Update() {
	for {
		x := <-ug.tick
		if ug.exit {
			break
		}
		for _, v := range ug.ug {
			v.Update(x)
		}
	}
	close(ug.tick)
	fmt.Println("Update Group Closed succesfully")
}

func (ug *UpdateGroup) Exit() {
	ug.exit = true
	ug.tick <- true
}

type updater struct {
	limit int
	ug    *UpdateGroup
}

func (u *updater) SetLimit(i int) {
	if i > u.limit {
		u.limit = i
	}
}

func (u *updater) GetUpdateGroup() *UpdateGroup {
	return u.ug
}

/*
strategy.New(strategy.Parameters{})

strategy.Long(strategy.size{account,10%},buy, startegy.ONCLOSE, reduceonly}

*/
