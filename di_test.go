package cqrs_dig

import (
	"context"
	"testing"

	cqrs "github.com/mitz-it/golang-cqrs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

type Response struct{}
type Event struct{}
type Command struct{}
type Query struct{}
type Service struct{}

type EventHandler1 struct {
	Service *Service
}

func (h *EventHandler1) Handle(ctx context.Context, event *Event) error {
	return nil
}

func NewEventHandler1(service *Service) cqrs.IEventHandler[*Event] {
	return &EventHandler1{
		Service: service,
	}
}

type EventHandler2 struct {
	Service *Service
}

func (h *EventHandler2) Handle(ctx context.Context, event *Event) error {
	return nil
}

func NewEventHandler2(service *Service) cqrs.IEventHandler[*Event] {
	return &EventHandler2{
		Service: service,
	}
}

type CommandHandler struct {
	Service *Service
}

func (h *CommandHandler) Handle(ctx context.Context, command *Command) (r *Response, err error) {
	return
}

func NewCommandHandler(service *Service) cqrs.ICommandHandler[*Command, *Response] {
	return &CommandHandler{
		Service: service,
	}
}

type QueryHandler struct {
	Service *Service
}

func (h *QueryHandler) Handle(ctx context.Context, command *Query) (r *Response, err error) {
	return
}

func NewQueryHandler(service *Service) cqrs.IQueryHandler[*Query, *Response] {
	return &QueryHandler{
		Service: service,
	}
}

type Behavior struct {
	Service *Service
}

func (b *Behavior) Handle(ctx context.Context, request interface{}, next cqrs.NextFunc) (interface{}, error) {
	return nil, nil
}

func NewBehavior(service *Service) *Behavior {
	return &Behavior{
		Service: service,
	}
}

func TestProvideEventSubscriber_WhenHandlerHasDependency_ShouldProvideHandlerWithDependency(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(func() *Service {
		return &Service{}
	})

	var handler *EventHandler1

	// act
	ProvideEventSubscriber[*Event](container, NewEventHandler1)
	err := container.Invoke(func(h cqrs.IEventHandler[*Event]) {
		handler = h.(*EventHandler1)
	})

	// assert
	assert.Nil(t, err)
	assert.Implements(t, (*cqrs.IEventHandler[*Event])(nil), handler)
	assert.NotNil(t, handler.Service)
}

func TestProvideEventSubscriber_WhenHandlerConstructorIsNil_ShouldReturnError(t *testing.T) {
	// arrange
	container := dig.New()

	// act
	err := ProvideEventSubscriber[*Event](container, nil)

	// assert
	assert.Error(t, err)
}

func TestProvideEventSubscribers_WhenHandlersHaveDependency_ShouldProvideHandlersWithDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(func() *Service {
		return &Service{}
	})

	var handlers []cqrs.IEventHandler[*Event]

	// act
	err := ProvideEventSubscribers[*Event](container, NewEventHandler1, NewEventHandler2)

	container.Invoke(func(params struct {
		dig.In

		Handlers []cqrs.IEventHandler[*Event] `group:"handlers"`
	}) {
		handlers = params.Handlers
	})

	// assert
	assert.Nil(t, err)
	assert.NotEmpty(t, handlers)
	assert.Contains(t, handlers, &EventHandler1{
		Service: &Service{},
	})
	assert.Contains(t, handlers, &EventHandler2{
		Service: &Service{},
	})
}

func TestProvideEventSubscribers_WhenAnyConstructorIsNil_ShouldReturnError(t *testing.T) {
	// arrange
	container := dig.New()

	// act
	err := ProvideEventSubscribers[*Event](container, NewEventHandler2, nil)

	// assert
	assert.Error(t, err)
}

func TestProvideCommandHandler_WhenHandlerHasDependency_ShouldProvideHandlerWithDependency(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(func() *Service {
		return &Service{}
	})
	var handler *CommandHandler
	// act
	err := ProvideCommandHandler[*Command, *Response](container, NewCommandHandler)

	container.Invoke(func(h cqrs.ICommandHandler[*Command, *Response]) {
		handler = h.(*CommandHandler)
	})

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.Service)
}

func TestProvideCommandHandler_WhenConstructorIsNil_ShouldReturnError(t *testing.T) {
	// arrange
	container := dig.New()
	// act
	err := ProvideCommandHandler[*Command, *Response](container, nil)

	// assert
	assert.Error(t, err)
}

func TestProvideQueryHandler_WhenHandlerHasDependency_ShouldProvideHandlerWithDependency(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(func() *Service {
		return &Service{}
	})
	var handler *QueryHandler
	// act
	err := ProvideQueryHandler[*Query, *Response](container, NewQueryHandler)

	container.Invoke(func(h cqrs.IQueryHandler[*Query, *Response]) {
		handler = h.(*QueryHandler)
	})

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.Service)
}

func TestProvideQueryHandler_WhenConstructorIsNil_ShouldReturnError(t *testing.T) {
	// arrange
	container := dig.New()
	// act
	err := ProvideQueryHandler[*Query, *Response](container, nil)

	// assert
	assert.Error(t, err)
}

func TestProvideCommandBehavior_WhenBehaviorHasDependency_ShouldProvideBehaviorWithDependency(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(func() *Service {
		return &Service{}
	})
	var behavior *Behavior

	// act
	err := ProvideCommandBehavior[*Behavior](container, 0, NewBehavior)

	container.Invoke(func(b *Behavior) {
		behavior = b
	})

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, behavior)
	assert.NotNil(t, behavior.Service)
}

func TestProvideCommandBehavior_WhenConstructorIsNil_ShouldReturnError(t *testing.T) {
	// arrange
	container := dig.New()
	// act
	err := ProvideCommandBehavior[*Behavior](container, 0, nil)

	// assert
	assert.Error(t, err)
}

func TestProvideQueryBehavior_WhenBehaviorHasDependency_ShouldProvideBehaviorWithDependency(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(func() *Service {
		return &Service{}
	})
	var behavior *Behavior

	// act
	err := ProvideQueryBehavior[*Behavior](container, 0, NewBehavior)

	container.Invoke(func(b *Behavior) {
		behavior = b
	})

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, behavior)
	assert.NotNil(t, behavior.Service)
}

func TestProvideQueryBehavior_WhenConstructorIsNil_ShouldReturnError(t *testing.T) {
	// arrange
	container := dig.New()
	// act
	err := ProvideQueryBehavior[*Behavior](container, 0, nil)

	// assert
	assert.Error(t, err)
}
