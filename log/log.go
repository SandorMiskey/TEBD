// region: packages

package log

import (
	"fmt"
)

// endregion: packages
// region: types

// endregion: types
// region: defaults

// endregion: defaults
// region: constructor

// func NewC(name string) *Logger {
func New(name string) {
	_, _, caller := Trace()
	fmt.Println(caller)
}

// endregion: constructor
// region: dumps

// func (c *Config) Sdump() string {
// 	return spew.Sdump(c)
// }

// func (c *Config) JSON(prefix string, indent string) ([]byte, error) {
// 	return json.MarshalIndent(c, prefix, indent)
// }

// endregion: dumps
