package sprig

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"net/url"
	"testing"
)

func TestBase64EncodeDecode(t *testing.T) {
	magicWord := "coffee"
	expect := base64.StdEncoding.EncodeToString([]byte(magicWord))

	if expect == magicWord {
		t.Fatal("Encoder doesn't work.")
	}

	tpl := `{{b64enc "coffee"}}`
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
	tpl = fmt.Sprintf("{{b64dec %q}}", expect)
	if err := runt(tpl, magicWord); err != nil {
		t.Error(err)
	}
}

func TestBase32EncodeDecode(t *testing.T) {
	magicWord := "coffee"
	expect := base32.StdEncoding.EncodeToString([]byte(magicWord))

	if expect == magicWord {
		t.Fatal("Encoder doesn't work.")
	}

	tpl := `{{b32enc "coffee"}}`
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
	tpl = fmt.Sprintf("{{b32dec %q}}", expect)
	if err := runt(tpl, magicWord); err != nil {
		t.Error(err)
	}
}

func TestURLEncodeDecode(t *testing.T) {
	magicWord := "cof&+/?fee"
	expect := url.QueryEscape(magicWord)

	tpl := fmt.Sprintf(`{{urlEnc "%s"}}`, magicWord)
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
	tpl = fmt.Sprintf("{{urlDec %q}}", expect)
	if err := runt(tpl, magicWord); err != nil {
		t.Error(err)
	}
}
