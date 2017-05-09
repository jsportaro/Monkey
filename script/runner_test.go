package script

import (
	"os"
	"testing"
)

func TestRunner(t *testing.T) {
	var fileName = "./test_script.monkey"
	Run(os.Stdout, []string{fileName})
}
