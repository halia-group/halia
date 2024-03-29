/*
 *
 *  * MIT License
 *  *
 *  * Copyright (c) [2021] [xialeistudio]
 *  *
 *  * Permission is hereby granted, free of charge, to any person obtaining a copy
 *  * of this software and associated documentation files (the "Software"), to deal
 *  * in the Software without restriction, including without limitation the rights
 *  * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  * copies of the Software, and to permit persons to whom the Software is
 *  * furnished to do so, subject to the following conditions:
 *  *
 *  * The above copyright notice and this permission notice shall be included in all
 *  * copies or substantial portions of the Software.
 *  *
 *  * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  * SOFTWARE.
 *
 */

package util

type AttributeMap interface {
	SetAttribute(key string, val interface{})
	GetAttribute(key string) interface{}
	HasAttribute(key string) bool
	GetStringAttribute(key string) string
	GetBoolAttribute(key string) bool
	GetIntAttribute(key string) int
}

type DefaultAttributeMap struct {
	attrs map[string]interface{}
}

func (m *DefaultAttributeMap) initialize() {
	if m.attrs == nil {
		m.attrs = make(map[string]interface{})
	}
}
func (m *DefaultAttributeMap) SetAttribute(key string, val interface{}) {
	m.initialize()
	m.attrs[key] = val
}

func (m *DefaultAttributeMap) GetAttribute(key string) interface{} {
	m.initialize()
	return m.attrs[key]
}

func (m *DefaultAttributeMap) HasAttribute(key string) bool {
	m.initialize()
	return m.attrs[key] != nil
}

func (m *DefaultAttributeMap) GetStringAttribute(key string) string {
	m.initialize()

	val := m.attrs[key]
	if val == nil {
		return ""
	}
	return val.(string)
}

func (m *DefaultAttributeMap) GetBoolAttribute(key string) bool {
	m.initialize()

	val := m.attrs[key]
	if val == nil {
		return false
	}
	return val.(bool)
}

func (m *DefaultAttributeMap) GetIntAttribute(key string) int {
	m.initialize()

	val := m.attrs[key]
	if val == nil {
		return 0
	}
	return val.(int)
}
