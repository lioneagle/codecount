package counter

import (
//"fmt"
)

type CodeCounterFactory struct {
	maps map[string]CodeCounter
}

func NewCodeCounterFactory() *CodeCounterFactory {
	factory := &CodeCounterFactory{}
	factory.maps = make(map[string]CodeCounter)
	factory.maps["go"] = NewGoCodeCounter()
	factory.maps["cpp"] = NewCppCodeCounter()
	factory.maps["c"] = NewCppCodeCounter()
	factory.maps["java"] = NewCppCodeCounter()
	factory.maps["erlang"] = NewErlangCodeCounter()
	return factory
}

func (factory *CodeCounterFactory) NewCounter(name string) (counter CodeCounter, ok bool) {
	if counter, ok = factory.maps[name]; ok {
		counter.Clear()
	}
	return counter, ok
}
