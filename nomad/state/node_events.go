package state

import (
	"github.com/hashicorp/nomad/nomad/stream"
	"github.com/hashicorp/nomad/nomad/structs"
)

const (
	NodeRegistrationTopic = "NodeRegistrationEvent"
)

type NodeRegistrationEvent struct {
	Event      *structs.NodeEvent
	NodeStatus string
	// Node       *structs.Node
}

func NodeRegisterEventFromChanges(tx ReadTxn, changes Changes) ([]stream.Event, error) {
	var event stream.Event
	for _, change := range changes.Changes {
		switch change.Table {
		case "nodes":
			after := change.After.(*structs.Node)

			event = stream.Event{
				Topic: NodeRegistrationTopic,
				Index: changes.Index,
				Key:   after.ID,
				Payload: &NodeRegistrationEvent{
					Event:      after.Events[len(after.Events)-1],
					NodeStatus: after.Status,
					// Node:       after,
				},
			}
		default:
			// not expecting  another table change for registering a node
		}
	}
	return []stream.Event{event}, nil
}

func NodeDeregisterEventFromChanges(tx ReadTxn, changes Changes) ([]stream.Event, error) {

	return []stream.Event{}, nil
}
