package tgx

import (
	"log"
	"sync"
)

type ParamBuilder struct {
	params map[string]interface{}
	mu     sync.Mutex
}

func NewParamBuilder() *ParamBuilder {
	return &ParamBuilder{params: make(map[string]interface{})}
}

func (pb *ParamBuilder) Add(key string, value interface{}) *ParamBuilder {
	log.Printf("Adding key: %s with value: %v", key, value)
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if value == nil {
		log.Printf("Skipping key: %s because value is nil", key)
		return pb
	}

	switch v := value.(type) {
	case string:
		if v != "" {
			pb.params[key] = v
		} else {
			log.Printf("Skipping key: %s because value is an empty string", key)
		}
	case *string:
		if v != nil && *v != "" {
			pb.params[key] = *v
		} else {
			log.Printf("Skipping key: %s because value is nil or points to an empty string", key)
		}
	case int, int64, float64, bool:
		pb.params[key] = v
	case *int:
		if v != nil {
			pb.params[key] = *v
		}
	case *int64:
		if v != nil {
			pb.params[key] = *v
		}
	case *bool:
		if v != nil {
			pb.params[key] = *v
		}
	case *float64:
		if v != nil {
			pb.params[key] = *v
		}
	default:
		// pb.params[key] = value
		log.Printf("Unsupported type for key: %s, skipping", key)
	}
	return pb
}

func (pb *ParamBuilder) Build() map[string]interface{} {
	return pb.params
}
