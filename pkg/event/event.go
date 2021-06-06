package event

import (
    "github.com/junglemc/JungleTree/internal/configuration"
    "log"
    . "reflect"
)

type Event interface{}

type Listener interface {
    OnEvent(event Event)
}

type listenerRegistry map[Type][]Listener

var listeners = make(listenerRegistry)

func Register(event Event, listener Listener) {
    // Debug logging, print event registration on function call
    if configuration.Config().DebugMode && configuration.Config().Verbose {
        log.Printf("Registering event listener: event=%s, listener=%s", TypeOf(event).Name(), TypeOf(listener).Name())
    }

    v := listeners[TypeOf(event)]
    if v == nil {
        v = make([]Listener, 0)
    }
    v = append(v, listener)
    listeners[TypeOf(event)] = v
}

func Trigger(event Event) {
    // Run on a separate goroutine to avoid hogging the spawning thread
    // TODO: Perhaps use channels to get the cancellable return result? Not yet implemented
    go func() {
        v := listeners[TypeOf(event)]
        if v == nil {
            return
        }

        for _, l := range v {
            // For long events, auto async it.
            // TODO: Thread pooling
            go l.OnEvent(event)
        }
    }()
}
