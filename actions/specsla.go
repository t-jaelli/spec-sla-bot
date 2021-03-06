package actions

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Azure/buffalo-azure/sdk/eventgrid"
	"github.com/Azure/spec-sla-bot/messages"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/google/go-github/github"
)

// SpecslaSubscriber gathers responses to all Requests sent to a particular endpoint.
type SpecslaSubscriber struct {
	eventgrid.Subscriber
}

// NewSpecslaSubscriber instantiates SpecslaSubscriber for use in a `buffalo.App`.
func NewSpecslaSubscriber(parent eventgrid.Subscriber) (created *SpecslaSubscriber) {
	dispatcher := eventgrid.NewTypeDispatchSubscriber(parent)

	created = &SpecslaSubscriber{
		Subscriber: dispatcher,
	}

	dispatcher.Bind("Github.PullRequestEvent", created.ReceivePullRequestEvent)
	dispatcher.Bind("Github.IssueCommentEvent", created.ReceiveIssueCommentEvent)
	dispatcher.Bind(eventgrid.EventTypeWildcard, created.ReceiveDefault)

	return
}

// ReceivePullRequestEvent will respond to an `eventgrid.Event` carrying a serialized `PullRequestEvent` as its payload.
func (s *SpecslaSubscriber) ReceivePullRequestEvent(c buffalo.Context, e eventgrid.Event) error {
	var payload github.PullRequestEvent

	if err := json.Unmarshal(e.Data, &payload); err != nil {
		return c.Error(http.StatusBadRequest, errors.New("unable to unmarshal request data"))
	}
	c.Logger().Debug("HERE")
	messages.CheckAcknowledgement(payload)

	// Replace the code below with your logic
	return c.Render(200, render.JSON(map[string]string{"message": "Hopefully this works"}))
}

func (s *SpecslaSubscriber) ReceiveIssueCommentEvent(c buffalo.Context, e eventgrid.Event) error {
	var payload github.IssueCommentEvent

	if err := json.Unmarshal(e.Data, &payload); err != nil {
		return c.Error(http.StatusBadRequest, errors.New("unable to unmarshal request data"))
	}
	c.Logger().Debug("Check acknowledgement of comment on PR")
	messages.CheckAcknowledgementComment(payload)

	// Replace the code below with your logic
	return c.Render(200, render.JSON(map[string]string{"message": "Hopefully this works"}))
}

func (s *SpecslaSubscriber) ReceiveDefault(c buffalo.Context, e eventgrid.Event) error {
	c.Logger().Debug(e)
	return c.Render(200, render.JSON(map[string]string{"message": "Hopefully this works"}))
}
