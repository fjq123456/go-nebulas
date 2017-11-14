// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//

package pod

import (
	"fmt"

	"github.com/nebulasio/go-nebulas/consensus"
	log "github.com/sirupsen/logrus"
)

// CommitState presents the prepare stage in pod
type CommitState struct {
	sm *consensus.StateMachine
}

// NewCommitState create a new prepare state
func NewCommitState(sm *consensus.StateMachine) *CommitState {
	return &CommitState{
		sm: sm,
	}
}

func (state *CommitState) String() string {
	return fmt.Sprintf("CommitState %p", state)
}

// Event handle event.
func (state *CommitState) Event(e consensus.Event) (bool, consensus.State) {
	return false, nil
}

// Enter called when transiting to this state.
func (state *CommitState) Enter(data interface{}) {
	log.Debug("CommitState enter.")
}

// Leave called when leaving this state.
func (state *CommitState) Leave(data interface{}) {
	log.Debug("CommitState leave.")
}
