// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package gw

type Error string

func (e Error) Error() string {
	return string(e)
}
