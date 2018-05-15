package rufus

import (
	"testing"
)

func TestApp(t *testing.T) {
	app := App{}

	if err := app.LoadConfigAndRouter(); err != nil {
		panic(err)
	}
}
