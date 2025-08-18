package sprig

import (
	"testing"
)

func TestMustEnv(t *testing.T) {
	tpl := `{{ mustEnv "INVALID" }}`
	if err := runt(tpl, "foo"); err == nil {
		t.Errorf("expected error, got: %v", err)
	}

	t.Setenv("TMP", "OK")

	tpl = `{{ mustEnv "TMP" }}`
	if err := runt(tpl, "OK"); err != nil {
		t.Error(err)
	}

}
