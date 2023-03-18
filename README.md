# Go - CQRS Dig

Methods to register [golang-cqrs](https://github.com/mitz-it/golang-cqrs) assets into [dig](https://github.com/uber-go/dig) container.


## Status

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mitz-it_golang-cqrs-dig&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=mitz-it_golang-cqrs-dig) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=mitz-it_golang-cqrs-dig&metric=coverage)](https://sonarcloud.io/summary/new_code?id=mitz-it_golang-cqrs-dig) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=mitz-it_golang-cqrs-dig&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=mitz-it_golang-cqrs-dig)

## Installing

```bash
go get -u github.com/mitz-it/golang-cqrs-dig
```

## Usage

This package provides methods to register CQRS implementations with dig-injected services. For all the following samples, let's assume that every CQRS implementation receives a `*PingPongService` into the constructor function. To understand how to create the CQRS implementation, please read the documentation available [here](https://github.com/mitz-it/golang-cqrs/blob/main/README.md).

Create the service to be injected:

```go
type PingPongService struct {
}

func (s *PingPongService) Play(ping string, pong string) string {
  return fmt.Sprintf("%s %s", ping, pong)
}

func NewPingPongService() *PingPongService {
  return &PingPongService{}
}
```

Create and configure the `*dig.Container` instance:

```go
container := dig.New()
container.Provide(NewPingPongService)
```

## Commands Usage

```go
// The command handler constructor func with injected *PongService
func NewPingCommandHandler(service *PongService) cqrs.ICommandHandler[*PingCommand, *PingResponse] {
  return &CommandHandler{
    service: service,
  }
}

// register the command handler constructor into the DI container
cqrs_dig.ProvideCommandHandler[*PingCommand, *PingResponse](container, NewPingCommandHandler)
```

## Queries Usage

```go
// The query handler constructor func with injected *PongService
func NewPingQueryHandler(service *PingPongService) cqrs.IQueryHandler[*PingQuery, *PingResponse] {
  return &PingQueryHandler{
    service: service,
  }
}

// register the query handler constructor into the DI container
cqrs_dig.ProvideQueryHandler[*PingQuery, *PingResponse](container, NewPingQueryHandler)
```

## Behaviors Usage

```go
// The behavior constructor func with injected *PongService
func NewPingBehavior(service *PingPongService) *PingBehavior {
  return &PingBehavior{
    service: service,
  }
}

// register the behavior constructor for commands into the DI container
cqrs_dig.ProvideCommandBehavior[*PingBehavior](container, 0, NewPingBehavior)

// register the behavior constructor for queries into the DI container
cqrs_dig.ProvideQueryBehavior[*PingBehavior](container, 0, NewPingBehavior)
```

## Events Usage

```go
// The event handler constructor func with injected *PongService
func NewPingEventHandler(service *PingPongService) cqrs.IEvenHandler[*PingEvent] {
  return &PingEventHandler{
    service: service,
  }
}

func NewPingEventHandler2(service *PingPongService) cqrs.IEvenHandler[*PingEvent] {
  return &PingEventHandler2{
    service: service,
  }
}

// register the event handler constructor into the DI container
cqrs_dig.ProvideEventSubscriber[*PingEvent](container, NewPingEventHandler)

// use this method if you have more than one handler per event type
cqrs_dig.ProvideEventSubscribers[*PingEvent](container, NewPingEventHandler, NewPingEventHandler2)
```
