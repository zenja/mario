package level

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/zenja/mario/event"
)

type FSM interface {
	GetCurrentState() string
	Receive(e event.Event, objToPass interface{}) error
	AddTransition(from, to string, effect func(interface{})) error
	AddTrigger(e event.Event, transitions []TransCondition) error
}

type TransCondition struct {
	From      string
	To        string
	Condition func(interface{}) bool
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// BaseFSM
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type BaseFSM struct {
	currState  string
	states     []string
	effectMap  map[string]func(interface{})
	triggerMap map[event.Event][]TransCondition
}

func NewBaseFSM(initState string, states []string) (*BaseFSM, error) {
	// state name cannot contains '=' for '>'
	var initStateValid bool
	for _, s := range states {
		if s == initState {
			initStateValid = true
		}
		if strings.Contains(s, "=") || strings.Contains(s, ">") {
			return nil, errors.Errorf("state name cannot contains '=' or '>' but is %s", s)
		}
	}
	if initStateValid == false {
		return nil, errors.Errorf("init state (%s) not found in possible states", initState)
	}
	return &BaseFSM{
		currState:  initState,
		states:     states,
		effectMap:  make(map[string]func(interface{})),
		triggerMap: make(map[event.Event][]TransCondition),
	}, nil
}

func (bm *BaseFSM) GetCurrentState() string {
	return bm.currState
}

func (bm *BaseFSM) Receive(e event.Event, objToPass interface{}) error {
	// check for possible transitions and invoke effect function if has any
	ts := bm.triggerMap[e]
	for i := range ts {
		transStr := bm.getTransStr(ts[i].From, ts[i].To)
		if ts[i].From != bm.currState {
			continue
		}

		effect, ok := bm.effectMap[transStr]
		if !ok {
			return errors.Errorf("Event trigger defined for %s but no effect function found", transStr)
		}

		if ts[i].Condition(objToPass) {
			// pass self to effect function
			effect(objToPass)
			// and THEN change state (notice the sequence!)
			bm.currState = ts[i].To
		}
	}
	return nil
}

func (bm *BaseFSM) AddTransition(from, to string, effect func(interface{})) error {
	// assure that from/to state exists in possible states
	fromValid := bm.isValidState(from)
	if fromValid == false {
		return errors.Errorf("from state (%s) not in possible states", from)
	}
	toValid := bm.isValidState(to)
	if toValid == false {
		return errors.Errorf("to state (%s) not in possible states", to)
	}

	bm.effectMap[bm.getTransStr(from, to)] = effect
	return nil
}

func (bm *BaseFSM) AddTrigger(e event.Event, ts []TransCondition) error {
	for i := range ts {
		if bm.isValidState(ts[i].From) == false {
			return errors.Errorf("%s is not a valid state", ts[i].From)
		}
		if bm.isValidState(ts[i].To) == false {
			return errors.Errorf("%s is not a valid state", ts[i].To)
		}
	}
	// assert that transitions cannot have duplicated "from" event,
	// because any state can have only
	bm.triggerMap[e] = ts
	return nil
}

func (bm *BaseFSM) getTransStr(from, to string) string {
	return fmt.Sprintf("%s=>%s", from, to)
}

func (bm *BaseFSM) getEffect(from, to string) (func(interface{}), bool) {
	f, ok := bm.effectMap[bm.getTransStr(from, to)]
	return f, ok
}

func (bm *BaseFSM) isValidState(s string) bool {
	var isValid bool
	for _, st := range bm.states {
		if s == st {
			isValid = true
			break
		}
	}
	return isValid
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// heroFSM
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type heroFSM struct {
	*BaseFSM
	hero *hero
}

// NewHeroFSM
func NewHeroFSM(hero *hero, initState string, states []string) (FSM, error) {
	bm, err := NewBaseFSM(initState, states)
	if err != nil {
		return nil, err
	}
	return &heroFSM{
		BaseFSM: bm,
		hero:    hero,
	}, nil
}
