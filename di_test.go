package cqrs_dig_test

import (
	"context"
	"fmt"
	"testing"

	cqrs_dig "github.com/mitz-it/golang-cqrs-dig"

	cqrs "github.com/mitz-it/golang-cqrs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

type PingPongService struct {
}

func (s *PingPongService) Play(ping string, pong string) string {
	return fmt.Sprintf("%s %s", ping, pong)
}

func NewPingPongService() *PingPongService {
	return &PingPongService{}
}

type PingResponse struct {
	Pong string
}

type PingCommand struct {
	Ping string
}

type PingQuery struct {
	Ping string
}

type PingEvent struct {
	Ping string
}

type PingCommandHandler struct {
	service *PingPongService
}

func (h *PingCommandHandler) Handle(ctx context.Context, command *PingCommand) (*PingResponse, error) {
	pong := h.service.Play(command.Ping, "pong")

	response := &PingResponse{
		Pong: pong,
	}

	return response, nil
}

func NewPingCommandHandler(service *PingPongService) cqrs.ICommandHandler[*PingCommand, *PingResponse] {
	return &PingCommandHandler{
		service: service,
	}
}

type PingQueryHandler struct {
	service *PingPongService
}

func (h *PingQueryHandler) Handle(ctx context.Context, query *PingQuery) (*PingResponse, error) {
	pong := h.service.Play(query.Ping, "pong")

	response := &PingResponse{
		Pong: pong,
	}

	return response, nil
}

func NewPingQueryHandler(service *PingPongService) cqrs.IQueryHandler[*PingQuery, *PingResponse] {
	return &PingQueryHandler{
		service: service,
	}
}

type PingBehavior struct {
	service *PingPongService
}

func (b *PingBehavior) Handle(ctx context.Context, request interface{}, next cqrs.NextFunc) (interface{}, error) {
	res, err := next()

	if err != nil {
		return nil, err
	}

	response := res.(*PingResponse)

	response.Pong = b.service.Play(response.Pong, "behavior also says pong")

	return response, nil
}

func NewPingBehavior(service *PingPongService) *PingBehavior {
	return &PingBehavior{
		service: service,
	}
}

type PingEventHandler struct {
	service *PingPongService
}

func (h *PingEventHandler) Handle(ctx context.Context, event *PingEvent) error {
	event.Ping = h.service.Play(event.Ping, "and event says pong!")
	return nil
}

func NewPingEventHandler(service *PingPongService) cqrs.IEventHandler[*PingEvent] {
	return &PingEventHandler{
		service: service,
	}
}

type PingEventHandler2 struct {
	service *PingPongService
}

func (h *PingEventHandler2) Handle(ctx context.Context, event *PingEvent) error {
	event.Ping = h.service.Play(event.Ping, "and event says pong!")
	return nil
}

func NewPingEventHandler2(service *PingPongService) cqrs.IEventHandler[*PingEvent] {
	return &PingEventHandler2{
		service: service,
	}
}

func Test_ProvideCommandHandler_WhenHasInjectedService_ShouldInvokeAllDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(NewPingPongService)

	command := &PingCommand{
		Ping: "ping",
	}

	// act
	cqrs_dig.ProvideCommandHandler[*PingCommand, *PingResponse](container, NewPingCommandHandler)

	response, _ := cqrs.Send[*PingCommand, *PingResponse](context.TODO(), command)

	// assert
	assert.Equal(t, "ping pong", response.Pong)
}

func Test_ProvideCommandHandler_WhenProvideCommandBehavior_ShouldInvokeAllDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(NewPingPongService)
	cqrs_dig.ProvideCommandHandler[*PingCommand, *PingResponse](container, NewPingCommandHandler)

	command := &PingCommand{
		Ping: "ping",
	}

	// act
	err := cqrs_dig.ProvideCommandBehavior[*PingBehavior](container, 0, NewPingBehavior)
	response, _ := cqrs.Send[*PingCommand, *PingResponse](context.TODO(), command)

	// assert
	assert.Equal(t, "ping pong behavior also says pong", response.Pong)
	assert.Nil(t, err)
}

func Test_ProvideQueryHandler_WhenHasInjectedService_ShouldInvokeAllDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(NewPingPongService)

	query := &PingQuery{
		Ping: "ping",
	}

	// act
	cqrs_dig.ProvideQueryHandler[*PingQuery, *PingResponse](container, NewPingQueryHandler)

	response, _ := cqrs.Request[*PingQuery, *PingResponse](context.TODO(), query)

	// assert
	assert.Equal(t, "ping pong", response.Pong)
}

func Test_ProvideQueryHandler_WhenProvideCommandBehavior_ShouldInvokeAllDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(NewPingPongService)
	cqrs_dig.ProvideQueryHandler[*PingQuery, *PingResponse](container, NewPingQueryHandler)

	query := &PingQuery{
		Ping: "ping",
	}

	// act
	err := cqrs_dig.ProvideQueryBehavior[*PingBehavior](container, 0, NewPingBehavior)
	response, _ := cqrs.Request[*PingQuery, *PingResponse](context.TODO(), query)

	// assert
	assert.Equal(t, "ping pong behavior also says pong", response.Pong)
	assert.Nil(t, err)
}

func Test_ProvideEventSubscriber_WhenHasInjectedDependencies_ShouldInvokeAllDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(NewPingPongService)

	event := &PingEvent{
		Ping: "ping",
	}

	// act
	cqrs_dig.ProvideEventSubscriber[*PingEvent](container, NewPingEventHandler)

	err := cqrs.PublishEvent(context.TODO(), event)

	// assert
	assert.Equal(t, "ping and event says pong!", event.Ping)
	assert.Nil(t, err)
}

func Test_ProvideEventSubscribers_WhenHasInjectedDependencies_ShouldInvokeAllDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(NewPingPongService)

	event := &PingEvent{
		Ping: "ping",
	}

	// act
	cqrs_dig.ProvideEventSubscribers[*PingEvent](container, NewPingEventHandler, NewPingEventHandler2)

	err := cqrs.PublishEvent(context.TODO(), event)

	// assert
	assert.Equal(t, "ping and event says pong! and event says pong!", event.Ping)
	assert.Nil(t, err)
}
