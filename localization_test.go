package sprig

import (
	"reflect"
	"testing"
)

func TestInitLocale(t *testing.T) {
	typeOfResult := reflect.TypeOf(InitLocale("EN")).Kind()
	ty := reflect.TypeOf(Localization{}).Kind()

	if typeOfResult != ty {
		t.Errorf("InitLocale function failed")
	}
}

func TestGetSetLocale(t *testing.T) {
	locale := InitLocale("EN")

	if locale.GetLocale() != "EN" {
		t.Errorf("GetLocale function failed")
	}

	locale.SetLocale("FR")

	if locale.GetLocale() != "FR" {
		t.Errorf("SetLocale function failed")
	}
}

func TestLocalization_SetTranslation(t *testing.T) {
	locale := InitLocale("EN")

	typeResult := reflect.TypeOf(locale.SetTranslation("EN", map[string]string{"test": "test"})).Kind()
	ty := reflect.TypeOf(Localization{}).Kind()

	if typeResult != ty {
		t.Errorf("SetTranslation function failed")
	}
}

func TestTranslateFunction(t *testing.T) {
	locale.SetTranslation("EN", map[string]string{"test": "en_test"})
	locale.SetTranslation("FR", map[string]string{"test": "fr_test"})
	tests := map[string]string{
		`{{ t "test" }}`:  "en_test",
		`{{ t "test2" }}`: "",
	}

	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}

	locale.SetLocale("FR")

	tests = map[string]string{
		`{{ t "test" }}`:  "fr_test",
		`{{ t "test2" }}`: "",
	}

	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}
