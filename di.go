package cqrs_dig

import (
	cqrs "github.com/mitz-it/golang-cqrs"
	"go.uber.org/dig"
)

type eventHandlerParams[TEvent any] struct {
	dig.In

	Handlers []cqrs.IEventHandler[TEvent] `group:"handlers"`
}

func ProvideEventSubscriber[TEvent any](container *dig.Container, constructor interface{}) error {
	err := container.Provide(constructor)

	if err != nil {
		return err
	}

	err = container.Invoke(func(handler cqrs.IEventHandler[TEvent]) error {
		return cqrs.RegisterEventSubscriber(handler)
	})

	return err
}

func ProvideEventSubscribers[TEvent any](container *dig.Container, constructors ...interface{}) error {
	for _, constructor := range constructors {
		err := container.Provide(constructor, dig.Group("handlers"))

		if err != nil {
			return err
		}
	}

	err := container.Invoke(func(params eventHandlerParams[TEvent]) error {
		return cqrs.RegisterEventSubscribers(params.Handlers...)
	})

	return err
}

func ProvideCommandHandler[TCommand any, TResponse any](container *dig.Container, constructor interface{}) error {
	err := container.Provide(constructor)

	if err != nil {
		return err
	}

	err = container.Invoke(func(handler cqrs.ICommandHandler[TCommand, TResponse]) error {
		return cqrs.RegisterCommandHandler(handler)
	})

	return err
}

func ProvideQueryHandler[TQuery any, TResponse any](container *dig.Container, constructor interface{}) error {
	err := container.Provide(constructor)

	if err != nil {
		return err
	}

	err = container.Invoke(func(handler cqrs.IQueryHandler[TQuery, TResponse]) error {
		return cqrs.RegisterQueryHandler(handler)
	})

	return err
}

func ProvideCommandBehavior[TBehavior cqrs.IBehavior](container *dig.Container, order int, constructor interface{}) error {
	err := container.Provide(constructor)

	if err != nil {
		return err
	}

	err = container.Invoke(func(behavior TBehavior) error {
		return cqrs.RegisterCommandBehavior(order, behavior)
	})

	return err
}

func ProvideQueryBehavior[TBehavior cqrs.IBehavior](container *dig.Container, order int, constructor interface{}) error {
	err := container.Provide(constructor)

	if err != nil {
		return err
	}

	err = container.Invoke(func(behavior TBehavior) error {
		return cqrs.RegisterQueryBehavior(order, behavior)
	})

	return err
}
