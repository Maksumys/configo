package configo

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseDefault(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		type confTest struct {
			Address string `configo:"address" default:"127.0.0.1"`
			Port    int    `configo:"port" default:"80"`
		}

		c, err := Parse[confTest](Option{})

		assert.NoError(t, err)
		assert.Equal(t, c.Address, "127.0.0.1")
		assert.Equal(t, c.Port, 80)
	})
	t.Run("struct into struct", func(t *testing.T) {
		type confTest struct {
			Http struct {
				Tls struct {
					Enabled bool `configo:"enabled" default:"true"`
				} `configo:"tls"`
				Address string `configo:"address" default:"127.0.0.1"`
				Port    int    `configo:"port" default:"80"`
			} `configo:"http"`
		}

		c, err := Parse[confTest](Option{})

		assert.NoError(t, err)
		assert.Equal(t, c.Http.Tls.Enabled, true)
		assert.Equal(t, c.Http.Address, "127.0.0.1")
		assert.Equal(t, c.Http.Port, 80)
	})
	t.Run("pointers", func(t *testing.T) {
		type confTest struct {
			Address *string `configo:"address" default:"127.0.0.1"`
			Port    *int    `configo:"port" default:"80"`
		}
		c, err := Parse[confTest](Option{})

		assert.NoError(t, err)
		assert.Equal(t, *c.Address, "127.0.0.1")
		assert.Equal(t, *c.Port, 80)
	})
	t.Run("struct pointers", func(t *testing.T) {
		type confTest struct {
			Http *struct {
				Tls struct {
					Enabled bool `configo:"enabled" default:"true"`
				} `configo:"tls"`
				Address string `configo:"address" default:"127.0.0.1"`
				Port    *int   `configo:"port" default:"80"`
			} `configo:"http"`
		}

		c, err := Parse[confTest](Option{})

		assert.NoError(t, err)
		assert.Equal(t, c.Http.Tls.Enabled, true)
		assert.Equal(t, c.Http.Address, "127.0.0.1")
		assert.Equal(t, *c.Http.Port, 80)
	})
}

func TestParseWithEnv(t *testing.T) {
	type confTest struct {
		Test struct {
			Env string `configo:"env" default:"test"`
		} `configo:"test"`
	}

	t.Run("with prefix", func(t *testing.T) {
		err := os.Setenv("CONFIGO_TEST2_TEST_ENV", "test_env")
		assert.NoError(t, err)

		c, err := Parse[confTest](Option{
			EnvPrefix:  "CONFIGO",
			EnvInclude: true,
			Key:        "test2",
		})
		assert.NoError(t, err)
		assert.Equal(t, c.Test.Env, "test_env")
	})
	t.Run("without prefix", func(t *testing.T) {
		err := os.Setenv("TEST_ENV", "test_env")
		assert.NoError(t, err)

		c, err := Parse[confTest](Option{
			EnvInclude: true,
		})

		assert.NoError(t, err)
		assert.Equal(t, c.Test.Env, "test_env")
	})
}
