package main

import (
	"context"
	"fmt"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"message,omitempty"`
}

type Receiver struct {
	client cloudevents.Client
	Target string `envconfig:"K_SINK"`
}

func main() {
	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatal(err.Error())
	}

	r := Receiver{client: client}
	if err := envconfig.Process("", &r); err != nil {
		log.Fatal(err.Error())
	}

	var receiver any
	if r.Target == "" {
		receiver = r.ReceiveAndReply
	} else {
		receiver = r.ReceiveAndSend
	}

	if err := client.StartReceiver(context.Background(), receiver); err != nil {
		log.Fatal(err)
	}
}

func newEvent() cloudevents.Event {
	r := cloudevents.NewEvent(cloudevents.VersionV1)
	r.SetID("42")
	r.SetType("app-event")
	r.SetSource("myapp")
	return r
}

func handle(req Request) Response {
	return Response{Message: fmt.Sprintf("Message received: %s", req.Message)}
}

func (recv *Receiver) ReceiveAndSend(ctx context.Context, event cloudevents.Event) cloudevents.Result {
	req := Request{}
	if err := event.DataAs(&req); err != nil {
		return cloudevents.NewHTTPResult(400, "failed to convert data: %s", err)
	}
	log.Printf("Got an event from: %q", req.Message)

	resp := handle(req)
	log.Printf("Sending event: %q", resp.Message)

	r := newEvent()
	if err := r.SetData("application/json", resp); err != nil {
		return cloudevents.NewHTTPResult(500, "failed to set response data: %s", err)
	}

	ctx = cloudevents.ContextWithTarget(ctx, recv.Target)
	return recv.client.Send(ctx, r)
}

func (recv *Receiver) ReceiveAndReply(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	req := Request{}
	if err := event.DataAs(&req); err != nil {
		return nil, cloudevents.NewHTTPResult(400, "failed to convert data: %s", err)
	}
	log.Printf("Got an event from: %q", req.Message)

	resp := handle(req)
	log.Printf("Replying with event: %q", resp.Message)

	r := newEvent()
	if err := r.SetData("application/json", resp); err != nil {
		return nil, cloudevents.NewHTTPResult(500, "failed to set response data: %s", err)
	}

	return &r, nil
}
