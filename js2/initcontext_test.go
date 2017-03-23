/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package js2

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/lib"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInitContextRequire(t *testing.T) {
	t.Run("k6", func(t *testing.T) {
		b, err := NewBundle(&lib.SourceData{
			Filename: "/script.js",
			Data: []byte(`
				import k6 from "k6";
				export let _k6 = k6;
				export let dummy = "abc123";
				export default function() {}
			`),
		}, afero.NewMemMapFs())
		if !assert.NoError(t, err, "bundle error") {
			return
		}

		rt, err := b.Instantiate()
		if !assert.NoError(t, err, "instance error") {
			return
		}
		assert.Contains(t, b.InitContext.Modules, "k6")

		exports := rt.Get("exports").ToObject(rt)
		if assert.NotNil(t, exports) {
			_, defaultOk := goja.AssertFunction(exports.Get("default"))
			assert.True(t, defaultOk, "default export is not a function")
			assert.Equal(t, "abc123", exports.Get("dummy").String())
		}

		k6 := rt.Get("_k6").ToObject(rt)
		if assert.NotNil(t, k6) {
			_, groupOk := goja.AssertFunction(k6.Get("group"))
			assert.True(t, groupOk, "k6.group is not a function")
		}

		t.Run("group", func(t *testing.T) {
			b, err := NewBundle(&lib.SourceData{
				Filename: "/script.js",
				Data: []byte(`
					import { group } from "k6";
					export let _group = group;
					export let dummy = "abc123";
					export default function() {}
				`),
			}, afero.NewMemMapFs())
			if !assert.NoError(t, err) {
				return
			}

			rt, err := b.Instantiate()
			if !assert.NoError(t, err) {
				return
			}
			assert.Contains(t, b.InitContext.Modules, "k6")

			exports := rt.Get("exports").ToObject(rt)
			if assert.NotNil(t, exports) {
				_, defaultOk := goja.AssertFunction(exports.Get("default"))
				assert.True(t, defaultOk, "default export is not a function")
				assert.Equal(t, "abc123", exports.Get("dummy").String())
			}

			_, groupOk := goja.AssertFunction(exports.Get("_group"))
			assert.True(t, groupOk, "{ group } is not a function")
		})
	})
}