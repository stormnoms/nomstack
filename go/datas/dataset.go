// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package datas

import (
	"regexp"

	"github.com/stormasm/noms/go/d"
	"github.com/stormasm/noms/go/types"
)

// DatasetRe is a regexp that matches a legal Dataset name anywhere within the
// target string.
var DatasetRe = regexp.MustCompile(`[a-zA-Z0-9\-_/]+`)

// DatasetFullRe is a regexp that matches a only a target string that is
// entirely legal Dataset name.
var DatasetFullRe = regexp.MustCompile("^" + DatasetRe.String() + "$")

// Dataset is a named Commit within a Database.
type Dataset struct {
	store   Database
	id      string
	headRef types.Ref
}

// Database returns the Database object in which this Dataset is stored.
// WARNING: This method is under consideration for deprecation.
func (ds Dataset) Database() Database {
	return ds.store
}

// ID returns the name of this Dataset.
func (ds Dataset) ID() string {
	return ds.id
}

// MaybeHead returns the current Head Commit of this Dataset, which contains
// the current root of the Dataset's value tree, if available. If not, it
// returns a new Commit and 'false'.
func (ds Dataset) MaybeHead() (types.Struct, bool) {
	if r, ok := ds.MaybeHeadRef(); ok {
		return r.TargetValue(ds.Database()).(types.Struct), true
	}
	return types.Struct{}, false
}

// Head returns the current head Commit, which contains the current root of
// the Dataset's value tree.
func (ds Dataset) Head() types.Struct {
	c, ok := ds.MaybeHead()
	d.PanicIfFalse(ok, "Dataset \"%s\" does not exist", ds.id)
	return c
}

// MaybeHeadRef returns the Ref of the current Head Commit of this Dataset,
// which contains the current root of the Dataset's value tree, if available.
// If not, it returns an empty Ref and 'false'.
func (ds Dataset) MaybeHeadRef() (types.Ref, bool) {
	return ds.headRef, ds.headRef != types.Ref{}
}

// HeadRef returns the Ref of the current head Commit, which contains the
// current root of the Dataset's value tree.
func (ds Dataset) HeadRef() types.Ref {
	r, ok := ds.MaybeHeadRef()
	d.PanicIfFalse(ok, "Dataset \"%s\" does not exist", ds.id)
	return r
}

// MaybeHeadValue returns the Value field of the current head Commit, if
// available. If not it returns nil and 'false'.
func (ds Dataset) MaybeHeadValue() (types.Value, bool) {
	if c, ok := ds.MaybeHead(); ok {
		return c.Get(ValueField), true
	}
	return nil, false
}

// HeadValue returns the Value field of the current head Commit.
func (ds Dataset) HeadValue() types.Value {
	c := ds.Head()
	return c.Get(ValueField)
}

func IsValidDatasetName(name string) bool {
	return DatasetFullRe.MatchString(name)
}
