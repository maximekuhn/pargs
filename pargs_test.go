package pargs_test

import (
	"errors"
	"testing"

	"github.com/maximekuhn/pargs"
)

func Test_ErrorInvalidTargetType(t *testing.T) {
	t.Run("target is a number", func(t *testing.T) {
		t.Parallel()
		err := pargs.Parse(42, nil)
		assertError(t, err)
		assertErrorIs(t, err, pargs.ErrInvalidTargetType)
	})

	t.Run("target is a pointer to a number", func(t *testing.T) {
		t.Parallel()
		x := 12
		err := pargs.Parse(&x, nil)
		assertError(t, err)
		assertErrorIs(t, err, pargs.ErrInvalidTargetType)
	})
}

func TestErrorStructFieldNotTagged(t *testing.T) {
	t.Run("no struct field tagged", func(t *testing.T) {
		t.Parallel()

		type args struct {
			Verbose    bool
			Quiet      bool
			ConfigPath string
		}

		var arg args
		err := pargs.Parse(&arg, nil)
		assertError(t, err)
		assertErrorIs(t, err, pargs.ErrStructFieldNotTagged)
	})
}

func Test_ItWorks(t *testing.T) {
	t.Run("mandatory boolean", func(t *testing.T) {
		t.Parallel()

		type args struct {
			Verbose bool `pargs:"flag:verbose"`
		}

		var arg args
		err := pargs.Parse(&arg, &pargs.Opts{
			Input: []string{"--verbose"},
		})
		assertNoError(t, err)
		assertTrue(t, arg.Verbose)
	})

	t.Run("short boolean flag", func(t *testing.T) {
		t.Parallel()

		type args struct {
			Verbose bool `pargs:"flag:v"`
		}

		var arg args
		err := pargs.Parse(&arg, &pargs.Opts{
			Input: []string{"-v"},
		})
		assertNoError(t, err)
		assertTrue(t, arg.Verbose)
	})

	t.Run("flag value different than field name", func(t *testing.T) {
		t.Parallel()

		type args struct {
			TotallyDifferentFieldName bool `pargs:"flag:default-config"`
		}

		var arg args
		err := pargs.Parse(&arg, &pargs.Opts{
			Input: []string{"--default-config"},
		})
		assertNoError(t, err)
		assertTrue(t, arg.TotallyDifferentFieldName)
	})
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func assertErrorIs(t *testing.T, got, want error) {
	if !errors.Is(got, want) {
		t.Fatalf("want '%v', got '%v'", want, got)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("expected nil error, got '%v'", err)
	}
}

func assertTrue(t *testing.T, b bool) {
	if !b {
		t.Fatal("expected true, got false")
	}
}
