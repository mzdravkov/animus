// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package be

import (
	// "fmt"

	. "github.com/gocircuit/escher/circuit"
)

// Required matter: Index, View, Verb
func materializeVerb(given Reflex, matter Circuit) (residue interface{}) {
	val, verb := lookup(matter)
	switch verb {
	case "*":
		return route(val, given, newSubMatter(matter))
	case "@":
		return materializeNoun(given, newSubMatter(matter).Grow("Noun", val))
	}
	panicWithMatter(matter, "unknown verb (%v)", verb)
	return
}

func newSubMatter(matter Circuit) Circuit {
	return New().
		Grow("Index", matter.CircuitAt("Index")).
		Grow("View", matter.CircuitAt("View")).
		Grow("Super", matter)
}

func relativize(matter Circuit) []Name {
	sup, ok := matter.CircuitOptionAt("Super")
	if !ok {
		return nil
	}
	if !sup.Has("Circuit") {
		return nil
	}
	supsup, ok := sup.CircuitOptionAt("Super")
	if !ok {
		return nil
	}
	supverb, ok := supsup.CircuitOptionAt("Resolved")
	if !ok {
		return nil
	}
	reladdr := Verb(supverb).Address()
	if len(reladdr) < 2 {
		return nil
	}
	return reladdr[:len(reladdr)-1] // chop off the circuit name at the end
}

func lookup(matter Circuit) (interface{}, string) {
	index, syntax := Index(matter.CircuitAt("Index")), matter.CircuitAt("Verb")
	verb, addr := Verb(syntax).Verb().(string), Verb(syntax).Address()

	rel := relativize(matter)
	var val interface{}
	if len(rel) > 0 {
		abs := append(rel, addr...)
		val = index.Recall(abs...) // lookup relative to enclosing circuit's parent circuit
		if val != nil {
			matter.Grow("Resolved", Circuit(NewVerbAddress(verb, abs...)))
			return val, verb
		}
	}
	val = index.Recall(addr...) // otherwise lookup globally
	matter.Include("Resolved", New().Grow(0, "???"))
	if val == nil {
		panicWithMatter(matter, "dangling address %v", Verb(syntax))
	}
	matter.Include("Resolved", Circuit(NewVerbAddress(verb, addr...)))
	return val, verb
}
