// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/attic-labs/testify/assert"
	"github.com/stormasm/noms/go/spec"
)

const (
	localSpec   = ldbSpec
	remoteSpec  = httpSpec
	testDs      = "testds"
	testObject  = "#pckdvpvr9br1fie6c3pjudrlthe7na18"
)

type testData struct {
	input    string
	expected string
}

var (
	rtestRoot = os.TempDir()

	rtestConfig = &Config{
		"",
		map[string]DbConfig{
			DefaultDbAlias: { localSpec },
			remoteAlias: { remoteSpec },
		},
	}

	dbTestsNoAliases = []testData {
		{localSpec, localSpec},
		{remoteSpec, remoteSpec},
	}

	dbTestsWithAliases = []testData {
		{"", localSpec},
		{remoteAlias, remoteSpec},
	}

	pathTestsNoAliases = []testData {
		{remoteSpec + "::" + testDs, remoteSpec + "::" + testDs},
		{remoteSpec + "::" + testObject, remoteSpec + "::" + testObject},

	}

	pathTestsWithAliases = []testData {
		{testDs, localSpec + "::" + testDs},
		{remoteAlias + "::" + testDs, remoteSpec + "::" + testDs},
		{testObject, localSpec + "::" + testObject},
		{remoteAlias + "::" + testObject, remoteSpec + "::" + testObject},
	}

)


func withConfig(t *testing.T) *Resolver {
	assert := assert.New(t)
	dir := filepath.Join(rtestRoot, "with-config")
	_, err := rtestConfig.WriteTo(dir)
	assert.NoError(err, dir)
	assert.NoError(os.Chdir(dir))
	r := NewResolver() // resolver must be created after changing directory
	return r

}

func withoutConfig(t *testing.T) *Resolver {
	assert := assert.New(t)
	dir := filepath.Join(rtestRoot, "without-config")
	assert.NoError(os.MkdirAll(dir, os.ModePerm), dir)
	assert.NoError(os.Chdir(dir))
	r := NewResolver() // resolver must be created after changing directory
	return r
}

func assertPathSpecsEquiv(assert *assert.Assertions, expected string, actual string) {
	e, err := spec.ParsePathSpec(expected)
	assert.NoError(err)
	a, err := spec.ParsePathSpec(actual)
	assert.NoError(err)
	assertDbSpecsEquiv(assert, e.DbSpec.String(), a.DbSpec.String())
	assert.Equal(e.Path.String(), a.Path.String())
}

func TestResolveDatabaseWithConfig(t *testing.T) {
	spec := withConfig(t)
	assert := assert.New(t)
	for _, d := range append(dbTestsNoAliases, dbTestsWithAliases...) {
		db := spec.ResolveDbSpec(d.input)
		assertDbSpecsEquiv(assert, d.expected, db)
	}
}

func TestResolvePathWithConfig(t *testing.T) {
	spec := withConfig(t)
	assert := assert.New(t)
	for _, d := range append(pathTestsNoAliases, pathTestsWithAliases...) {
		path := spec.ResolvePathSpec(d.input)
		assertPathSpecsEquiv(assert, d.expected, path)
	}
}

func TestResolveDatabaseWithoutConfig(t *testing.T) {
	spec := withoutConfig(t)
	assert := assert.New(t)
	for _, d := range dbTestsNoAliases {
		db := spec.ResolveDbSpec(d.input)
		assert.Equal(d.expected, db, d.input)
	}
}

func TestResolvePathWithoutConfig(t *testing.T) {
	spec := withoutConfig(t)
	assert := assert.New(t)
	for _, d := range pathTestsNoAliases {
		path := spec.ResolvePathSpec(d.input)
		assertPathSpecsEquiv(assert, d.expected, path)
	}

}

func TestResolveDestPathWithDot(t *testing.T) {
	spec := withConfig(t)
	assert := assert.New(t)

	data := []struct {
		src string
		dest string
		expSrc string
		expDest string
	} {
		{testDs, remoteSpec+"::.", 	localSpec+"::"+testDs, remoteSpec+"::"+testDs},
		{remoteSpec+"::"+testDs, ".",	remoteSpec+"::"+testDs, localSpec+"::"+testDs},
	}
	for _, d := range data {
		src := spec.ResolvePathSpec(d.src)
		dest := spec.ResolvePathSpec(d.dest)
		assertPathSpecsEquiv(assert, d.expSrc, src)
		assertPathSpecsEquiv(assert, d.expDest, dest)
	}

}
