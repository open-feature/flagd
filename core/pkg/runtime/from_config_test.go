package runtime

import (
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/require"
)

func Test_setupJSONEvaluator(t *testing.T) {
	lg := logger.NewLogger(nil, false)

	je := setupJSONEvaluator(lg, store.NewFlags())
	require.NotNil(t, je)
}
