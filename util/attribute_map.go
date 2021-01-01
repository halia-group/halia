package util

type AttributeMap interface {
	SetAttribute(key string, val interface{})
	GetAttribute(key string) interface{}
	HasAttribute(key string) bool
	GetStringAttribute(key string) string
	GetBoolAttribute(key string) bool
}

type DefaultAttributeMap struct {
	attrs map[string]interface{}
}

func NewDefaultAttributeMap() *DefaultAttributeMap {
	return &DefaultAttributeMap{attrs: make(map[string]interface{})}
}

func (m *DefaultAttributeMap) SetAttribute(key string, val interface{}) {
	m.attrs[key] = val
}

func (m *DefaultAttributeMap) GetAttribute(key string) interface{} {
	return m.attrs[key]
}

func (m *DefaultAttributeMap) HasAttribute(key string) bool {
	return m.attrs[key] != nil
}

func (m *DefaultAttributeMap) GetStringAttribute(key string) string {
	val := m.attrs[key]
	if val == nil {
		return ""
	}
	return val.(string)
}

func (m *DefaultAttributeMap) GetBoolAttribute(key string) bool {
	val := m.attrs[key]
	if val == nil {
		return false
	}
	return val.(bool)
}
