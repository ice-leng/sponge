package sse

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/go-dev-frame/sponge/pkg/utils"
)

func TestSSEClient_Connect(t *testing.T) {
	port, _ := utils.GetAvailablePort()
	eventType := DefaultEventType
	ctx, cancel := context.WithCancel(context.Background())

	// start sse server
	hub := NewHub(WithContext(ctx, cancel))
	go runSSEServer(port, hub)
	defer hub.Close()

	time.Sleep(100 * time.Millisecond) // wait for server to start

	// create sse client
	client := NewClient(fmt.Sprintf("http://localhost:%d/events", port),
		WithClientLogger(zap.NewExample()),
		WithClientHeaders(map[string]string{"Authorization": "Bearer abcdef"}),
		WithClientReconnectTimeInterval(time.Millisecond*100),
	)
	var receivedEvent *Event
	client.OnEvent(eventType, func(event *Event) {
		t.Log("on event", event)
		receivedEvent = event
	})
	err := client.Connect()
	assert.NoError(t, err)
	defer client.Disconnect()

	time.Sleep(100 * time.Millisecond)

	// push event to sse client
	var event = &Event{
		Event: eventType,
		Data:  "test-push-data",
	}
	err = hub.Push(nil, event)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond) // wait for event to be received

	assert.Equal(t, event.Event, receivedEvent.Event)
}

func TestSSEClient_Reconnect(t *testing.T) {
	port, _ := utils.GetAvailablePort()
	eventType := DefaultEventType
	ctx, cancel := context.WithCancel(context.Background())
	hub := NewHub(WithContext(ctx, cancel))
	go runSSEServer2(ctx, port, hub)

	client := NewClient(fmt.Sprintf("http://localhost:%d/events", port), WithClientReconnectTimeInterval(time.Millisecond*100))
	client.OnEvent(eventType, func(event *Event) {
		t.Log("on event", event)
	})
	err := client.Connect()
	assert.NoError(t, err)
	time.Sleep(200 * time.Millisecond)

	// expected connect is true
	assert.True(t, client.GetConnectStatus(), "Client should be connected")

	// close sse server
	hub.Close()

	time.Sleep(time.Millisecond * 300)
	// expected connect is false
	assert.False(t, client.GetConnectStatus(), "Client should be disconnected")

	// run sse server again
	ctx, cancel = context.WithCancel(context.Background())
	hub = NewHub(WithContext(ctx, cancel))
	go runSSEServer2(ctx, port, hub)
	defer hub.Close()
	_ = client.Connect()

	// wait for client to reconnect
	time.Sleep(2 * time.Second)

	// expected connect is true
	assert.True(t, client.GetConnectStatus(), "Client should be connected again")
	time.Sleep(200 * time.Millisecond)

	client.Disconnect()
}
