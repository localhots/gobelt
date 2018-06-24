package filecache

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Load will look for a given file. If it exists, it would unmarshal its
// contents into the given value. If the file does not exist, a supplied
// function would be called and its output would be written to a file.
func Load(val interface{}, filename string, fn func()) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fn()
		body, err := json.Marshal(val)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(filename, body, 0766)
	}
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, val)
}
