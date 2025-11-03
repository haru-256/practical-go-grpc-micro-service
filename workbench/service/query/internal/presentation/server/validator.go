package server

import "buf.build/go/protovalidate"

// TODO: interceptorを実装
func NewValidator() (protovalidate.Validator, error) {
	v, err := protovalidate.New()
	if err != nil {
		return nil, err
	}
	return v, nil
}
