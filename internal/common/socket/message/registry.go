package message

import "reflect"

// Registry 消息类型注册
type Registry struct {
	typeMap map[Type]reflect.Type
}

// NewRegistry 新建实例
func NewRegistry() *Registry {
	r := &Registry{
		typeMap: make(map[Type]reflect.Type),
	}

	return r
}

// Add 新增
func (r *Registry) Add(msgType Type, payload interface{}) {
	r.typeMap[msgType] = reflect.TypeOf(payload).Elem()
}

// New 根据类型创建实例
func (r *Registry) New(msgType Type) (v interface{}, found bool) {
	elem, ok := r.typeMap[msgType]
	if !ok {
		return nil, false
	}

	return reflect.New(elem).Interface(), true
}
