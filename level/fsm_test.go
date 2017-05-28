package level_test

import (
	"testing"

	"github.com/zenja/mario/level"
)

func TestBaseFSMHappyPath(t *testing.T) {
	bm, err := level.NewBaseFSM("init", []string{"init", "running", "dead"})
	assert(err, "failed to create BaseFSM instance", t)

	var actual string

	var eventLeft = "event-left"
	var eventRight = "event-right"
	var eventSpace = "event-space"

	// add transitions
	err = bm.AddTransition("init", "running", func(_ interface{}) {
		actual = "I->R"
	})
	assert(err, "failed to add transition", t)
	err = bm.AddTransition("running", "dead", func(_ interface{}) {
		actual = "R->D"
	})
	assert(err, "failed to add transition", t)

	// add triggers
	err = bm.AddTrigger(eventLeft, []level.TransCondition{
		{"init", "running", func(o interface{}) bool {
			if o == "this-will-cause-condition-to-false" {
				return false
			} else {
				return true
			}
		}},
	})
	assert(err, "tailed to add trigger", t)
	err = bm.AddTrigger(eventRight, []level.TransCondition{
		{"running", "dead", func(_ interface{}) bool { return true }},
	})
	assert(err, "tailed to add trigger", t)

	err = bm.Receive(eventSpace, bm)
	assert(err, "error after receive SPACE event", t)
	if actual != "" {
		t.Fatalf("event not related to any trigger should not have effect")
	}

	err = bm.Receive(eventRight, bm)
	assert(err, "error after receive RIGHT event", t)
	if actual != "" {
		t.Fatalf("event on state which is not related to any trigger for thie event should not have effect")
	}
	if bm.GetCurrentState() != "init" {
		t.Fatalf("expected state: %s; actual: %s", "init", bm.GetCurrentState())
	}

	err = bm.Receive(eventLeft, "this-will-cause-condition-to-false")
	assert(err, "error after receive LEFT event", t)
	if actual != "" {
		t.Fatalf("event on correct state but with false condition should not be triggered")
	}
	if bm.GetCurrentState() != "init" {
		t.Fatalf("expected state: %s; actual: %s", "init", bm.GetCurrentState())
	}

	err = bm.Receive(eventLeft, bm)
	assert(err, "error after receive LEFT event", t)
	if actual != "I->R" {
		t.Fatalf("init -> running transition should be triggered but not (expected: %s, actual: %s)", "I->R", actual)
	}
	if bm.GetCurrentState() != "running" {
		t.Fatalf("expected state: %s; actual: %s", "running", bm.GetCurrentState())
	}

	err = bm.Receive(eventRight, bm)
	assert(err, "error after receive RIGHT event", t)
	if actual != "R->D" {
		t.Fatalf("running -> dead transition should be triggered but not (expected: %s, actual: %s)", "R->D", actual)
	}
	if bm.GetCurrentState() != "dead" {
		t.Fatalf("expected state: %s; actual: %s", "dead", bm.GetCurrentState())
	}
}

func assert(err error, msg string, t *testing.T) {
	if err != nil {
		t.Fatalf("%s: %s", msg, err)
	}
}
