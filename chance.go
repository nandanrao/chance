package chance

import "sync"

func Pool(numWorkers int, tasks <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	chans := []<-chan interface{}{}
	for i := 1; i <= numWorkers; i++ {
		ch := Worker(tasks, fn)
		chans = append(chans, ch)
	}

	return Merge(chans...)
}

func Worker(inch <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func(){
		defer close(ch)
		for d := range inch {
			ch <- fn(d)
		}
	}()
	return ch
}

func Merge(cs ...<-chan interface{}) <-chan interface{} {
    var wg sync.WaitGroup
    out := make(chan interface{})

    output := func(c <-chan interface{}) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

func Channelify(items []interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for _, x := range items {
			ch <- x
		}
	}()
	return ch
}

func Chunk(ch chan interface{}, size int) chan []interface{} {
	var b []interface{}
	out := make(chan []interface{})

	go func(){
		for v := range ch {
			b = append(b, v)
			if len(b) == size {
				out <- b
				b = make([]interface{}, 0)
			}
		}
		// send the remaining partial buffer
		if len(b) > 0 {
			out <- b
		}
		close(out)
	}()

	return out
}
