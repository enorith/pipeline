package define

import "fmt"

type UndefinedTypeError string

func (e UndefinedTypeError) Error() string {
	return fmt.Sprintf("undefined type: %s", string(e))
}
