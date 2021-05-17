/*

 MIT License

 Copyright (c) 2021 Looker Data Sciences, Inc.

 Permission is hereby granted, free of charge, to any person obtaining a copy
 of this software and associated documentation files (the "Software"), to deal
 in the Software without restriction, including without limitation the rights
 to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 copies of the Software, and to permit persons to whom the Software is
 furnished to do so, subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 SOFTWARE.

 */

import { TestConfig } from './testUtils'
import { IEnumType } from './sdkModels'
import { GoGen } from './go.gen'

const config = TestConfig()
const apiTestModel = config.apiTestModel
const gen = new GoGen(apiTestModel)
const indent = ''

describe('Go generator', () => {
  describe('comment header', () => {
    it('is empty with no comment', () => {
      expect(gen.commentHeader(indent, '')).toEqual('')
    })

    it('is multiple lines with a two line comment', () => {
      const expected = `/*

foo
bar
*/
`
      expect(gen.commentHeader(indent, 'foo\nbar')).toEqual(expected)
    })
  })

  describe('types', () => {
    it('enum type', () => {
      const type = apiTestModel.types.PermissionType as IEnumType
      expect(type).toBeDefined()
      expect(type.values).toEqual(['view', 'edit'])
      const expected = `type PermissionType string
const PermissionType_View PermissionType = "view"
const PermissionType_Edit PermissionType = "edit"
`
      const actual = gen.declareType('', type)
      expect(actual).toEqual(expected)
    })

    it('special needs', () => {
      const type = apiTestModel.types.HyphenType
      const expected = `
type HyphenType struct {
  ProjectName     *string  \`json:"project_name,omitempty"\`      // A normal variable name
  ProjectDigest   *string  \`json:"project_digest,omitempty"\`    // A hyphenated property name
  ComputationTime *float32 \`json:"computation_time,omitempty"\`  // A spaced out property name
}`
      const actual = gen.declareType('', type)
      expect(actual).toEqual(expected)
    })
    it('noComment special needs', () => {
      const type = apiTestModel.types.HyphenType
      const expected = `
type HyphenType struct {
  ProjectName     *string  \`json:"project_name,omitempty"\`
  ProjectDigest   *string  \`json:"project_digest,omitempty"\`
  ComputationTime *float32 \`json:"computation_time,omitempty"\`
}`
      gen.noComment = true
      const actual = gen.declareType('', type)
      gen.noComment = false
      expect(actual).toEqual(expected)
    })
  })
})
