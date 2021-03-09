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
