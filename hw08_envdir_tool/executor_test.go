package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testEnviron = []string{
	"FOOBAR=123",
	"JJJ=jjj",
	"POO=kkk",
	"DDD=ddd",
}

func TestUpdateEnviron(t *testing.T) {
	t.Run("delete var", func(t *testing.T) {
		upd := Environment{
			"POO": {"", true},
		}

		env := UpdateEnviron(testEnviron, upd)

		require.NotContains(t, env, "POO=kkk")
		require.Len(t, env, len(testEnviron)-1)
	})

	t.Run("update with value", func(t *testing.T) {
		upd := Environment{
			"POO": {"jjj", false},
		}

		env := UpdateEnviron(testEnviron, upd)

		require.Contains(t, env, "POO=jjj")
		require.Len(t, env, len(testEnviron))
	})

	t.Run("add env var", func(t *testing.T) {
		upd := Environment{
			"KKK": {"111", false},
		}

		env := UpdateEnviron(testEnviron, upd)

		require.Contains(t, env, "KKK=111")
		require.Len(t, env, len(testEnviron)+1)
	})
}
