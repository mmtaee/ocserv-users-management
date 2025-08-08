package processor

import (
	"context"
)

func Processor(ctx context.Context, streamText <-chan string, broadcaster chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-streamText:
			if !ok {
				return
			}
			// SSE clients
			select {
			case broadcaster <- line:
			case <-ctx.Done():
				return
			}
		}
	}
}
