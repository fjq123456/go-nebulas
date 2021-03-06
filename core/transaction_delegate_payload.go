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

package core

import (
	"encoding/json"

	"github.com/nebulasio/go-nebulas/storage"
	"github.com/nebulasio/go-nebulas/util"
)

// Action Constants
const (
	DelegateAction   = "do"
	UnDelegateAction = "undo"
)

// DelegatePayload carry election information
type DelegatePayload struct {
	Action    string
	Delegatee string
}

var (
	vote = []byte("vote")
)

// LoadDelegatePayload from bytes
func LoadDelegatePayload(bytes []byte) (*DelegatePayload, error) {
	payload := &DelegatePayload{}
	if err := json.Unmarshal(bytes, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// NewDelegatePayload with function & args
func NewDelegatePayload(action string, addr string) *DelegatePayload {
	return &DelegatePayload{
		Action:    action,
		Delegatee: addr,
	}
}

// ToBytes serialize payload
func (payload *DelegatePayload) ToBytes() ([]byte, error) {
	return json.Marshal(payload)
}

// Execute the call payload in tx, call a function
func (payload *DelegatePayload) Execute(tx *Transaction, block *Block) (*util.Uint128, error) {
	delegator := tx.from.Bytes()
	counter := util.NewUint128()
	delegatee, err := AddressParse(payload.Delegatee)
	if err != nil {
		return counter, err
	}
	counter.Add(counter.Int, one.Int)
	// check delegatee valid
	_, err = block.dposContext.candidateTrie.Get(delegatee.Bytes())
	if err != nil && err != storage.ErrKeyNotFound {
		return counter, err
	}
	if err == storage.ErrKeyNotFound {
		return counter, ErrInvalidDelegateToNonCandidate
	}
	counter.Add(counter.Int, one.Int)
	pre, err := block.dposContext.voteTrie.Get(delegator)
	if err != nil && err != storage.ErrKeyNotFound {
		return counter, err
	}
	switch payload.Action {
	case DelegateAction:
		if err != storage.ErrKeyNotFound {
			key := append(pre, delegator...)
			counter.Add(counter.Int, one.Int)
			if _, err = block.dposContext.delegateTrie.Del(key); err != nil {
				return counter, err
			}
		}
		key := append(delegatee.Bytes(), delegator...)
		counter.Add(counter.Int, one.Int)
		if _, err = block.dposContext.delegateTrie.Put(key, delegator); err != nil {
			return counter, err
		}
		counter.Add(counter.Int, one.Int)
		if _, err = block.dposContext.voteTrie.Put(delegator, delegatee.Bytes()); err != nil {
			return counter, err
		}
	case UnDelegateAction:
		if !delegatee.address.Equals(pre) {
			return counter, ErrInvalidUnDelegateFromNonDelegatee
		}
		key := append(delegatee.Bytes(), delegator...)
		counter.Add(counter.Int, one.Int)
		if _, err = block.dposContext.delegateTrie.Del(key); err != nil {
			return counter, err
		}
		counter.Add(counter.Int, one.Int)
		if _, err = block.dposContext.voteTrie.Del(delegator); err != nil {
			return counter, err
		}
	default:
		return counter, ErrInvalidDelegatePayloadAction
	}
	return counter, nil
}
