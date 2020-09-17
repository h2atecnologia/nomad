package state

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/nomad/nomad/mock"
	"github.com/hashicorp/nomad/nomad/stream"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/stretchr/testify/require"
)

func TestNodeRegisterEventFromChanges(t *testing.T) {
	cases := []struct {
		Name       string
		MsgType    structs.MessageType
		Setup      func(s *StateStore, tx *txn) error
		Mutate     func(s *StateStore, tx *txn) error
		WantEvents []stream.Event
		WantErr    bool
		WantTopic  string
	}{
		{
			MsgType:   structs.NodeRegisterRequestType,
			WantTopic: NodeRegistrationTopic,
			Name:      "node registered",
			Mutate: func(s *StateStore, tx *txn) error {
				return upsertNodeTxn(tx, tx.Index, mock.Node())
			},
			WantEvents: nil,
			WantErr:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			s := TestStateStoreCfg(t, TestStateStorePublisher(t))
			defer s.StopEventPublisher()

			if tc.Setup != nil {
				// Bypass publish mechanism for setup
				setupTx := s.db.WriteTxn(10)
				require.NoError(t, tc.Setup(s, setupTx))
				setupTx.Txn.Commit()
			}

			tx := s.db.WriteTxn(100)
			require.NoError(t, tc.Mutate(s, tx))

			changes := Changes{Changes: tx.Changes(), Index: 100, MsgType: tc.MsgType}
			got, err := processDBChanges(tx, changes)
			spew.Dump(got)

			if tc.WantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tc.WantEvents, got)

			for _, g := range got {
				require.Equal(t, tc.WantTopic, g.Topic)
				require.Equal(t, 100, g.Index)
			}
		})
	}
}

func TestNodeDeregisterEventFromChanges(t *testing.T) {

}
