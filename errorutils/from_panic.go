package errorutils

import "github.com/rotisserie/eris"

func FromPanic(v any) error {
	if v == nil {
		return nil
	}

	if err, ok := v.(error); ok {
		return err
	}
	
	return eris.Errorf("recovered from panic: %v", v)
}