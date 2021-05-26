package predict

import "testing"

import "github.com/stretchr/testify/assert"

func TestConfigCheck(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		cfg := Options(OptValues("foo", "bar", "foo-bar"), OptCheck())
		assert.NoError(t, cfg.Check("foo"))
		assert.NoError(t, cfg.Check("bar"))
		assert.NoError(t, cfg.Check("foo-bar"))
		assert.Error(t, cfg.Check("fo"))
		assert.Error(t, cfg.Check("baz"))
	})

	t.Run("disabled", func(t *testing.T) {
		cfg := Options(OptValues("foo", "bar", "foo-bar"))
		assert.NoError(t, cfg.Check("foo"))
		assert.NoError(t, cfg.Check("fo"))
		assert.NoError(t, cfg.Check("baz"))
	})
}

func TestConfigPredict(t *testing.T) {
	t.Parallel()

	t.Run("set", func(t *testing.T) {
		cfg := Options(OptValues("foo", "bar", "foo-bar"))
		assert.Equal(t, []string{"foo", "bar", "foo-bar"}, cfg.Predict(""))
	})

	t.Run("not set", func(t *testing.T) {
		cfg := Options()
		assert.Nil(t, cfg.Predict(""))
	})
}
