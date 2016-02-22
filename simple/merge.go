package simple

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "sync"

func MergeBool(cs ...<-chan bool) <-chan bool {
	var wg sync.WaitGroup
	out := make(chan bool)
	output := func(c <-chan bool) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	go func() { wg.Wait(); close(out) }()
	return out
}
