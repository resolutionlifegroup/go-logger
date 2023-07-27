package log

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestMarshalStackForGCP(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		stack := marshalStackForGCP(nil)

		require.Nil(t, stack)
	})

	t.Run("new error", func(t *testing.T) {
		err := errors.New("new test error")
		stack := marshalStackForGCP(err)

		require.NotNil(t, stack)
		require.NotEmpty(t, stack)

		stackS, isString := stack.(string)
		require.True(t, isString, "the returned value is not a string")
		require.True(t, strings.HasPrefix(stackS, "goroutine 1 [running]:"))
	})
}
