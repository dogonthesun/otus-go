package hw06pipelineexecution

type (
	In  = <-chan any
	Out = In
	Bi  = chan any
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, cancel In, stages ...Stage) Out {
	wrapper := func(in In, cancel In) Out {
		out := make(Bi)
		go func() {
			defer func() {
				close(out)
				for range in { //nolint:all
				}
			}()
			for {
				select {
				case <-cancel:
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					select {
					case out <- v:
					case <-cancel:
						return
					}
				}
			}
		}()
		return out
	}

	out := in
	for _, stage := range stages {
		out = stage(wrapper(out, cancel))
	}

	return wrapper(out, cancel)
}
