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

// BaseGasCount returns base gas count
func (payload *DelegatePayload) BaseGasCount() *util.Uint128 {
	return DelegateBaseGasCount
}

// Execute the call payload in tx, call a function
func (payload *DelegatePayload) Execute(ctx *PayloadContext) (*util.Uint128, error) {
	delegator := ctx.tx.from.Bytes()
	delegatee, err := AddressParse(payload.Delegatee)
	if err != nil {
		return ZeroGasCount, err
	}
	// check delegatee valid
	_, err = ctx.block.dposContext.candidateTrie.Get(delegatee.Bytes())
	if err != nil && err != storage.ErrKeyNotFound {
		return ZeroGasCount, err
	}
	if err == storage.ErrKeyNotFound {
		return ZeroGasCount, ErrInvalidDelegateToNonCandidate
	}
	pre, err := ctx.block.dposContext.voteTrie.Get(delegator)
	if err != nil && err != storage.ErrKeyNotFound {
		return ZeroGasCount, err
	}
	switch payload.Action {
	case DelegateAction:
		if err != storage.ErrKeyNotFound {
			key := append(pre, delegator...)
			if _, err = ctx.block.dposContext.delegateTrie.Del(key); err != nil {
				return ZeroGasCount, err
			}
		}
		key := append(delegatee.Bytes(), delegator...)
		if _, err = ctx.block.dposContext.delegateTrie.Put(key, delegator); err != nil {
			return ZeroGasCount, err
		}
		if _, err = ctx.block.dposContext.voteTrie.Put(delegator, delegatee.Bytes()); err != nil {
			return ZeroGasCount, err
		}
	case UnDelegateAction:
		if !delegatee.address.Equals(pre) {
			return ZeroGasCount, ErrInvalidUnDelegateFromNonDelegatee
		}
		key := append(delegatee.Bytes(), delegator...)
		if _, err = ctx.block.dposContext.delegateTrie.Del(key); err != nil {
			return ZeroGasCount, err
		}
		if _, err = ctx.block.dposContext.voteTrie.Del(delegator); err != nil {
			return ZeroGasCount, err
		}
	default:
		return ZeroGasCount, ErrInvalidDelegatePayloadAction
	}
	return ZeroGasCount, nil
}
