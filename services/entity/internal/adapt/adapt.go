// Package adapt turns one input message into a normalised Observation (fusion-and-tracking.md
// Stage 0). Each adapter returns (observation, "", nil), (nil, skipReason, nil) for a message
// the service ignores on purpose, or (nil, "", err) for an invalid message, which the caller
// dead-letters with the reason.
package adapt

import (
	"strings"

	"github.com/onemindos/omos/services/entity/internal/model"
)

// Skip reasons, counted by the service.
const (
	SkipEcho        = "echo"         // the projector's own CoT
	SkipUnsupported = "unsupported"  // a subject or CoT type this slice does not read
	SkipNoPosition  = "no_position"  // an intel record without coordinates
	SkipOutOfScope  = "out_of_scope" // military surge subjects (BUS-DESIGN D9)
)

// Parse dispatches on the subject.
func Parse(subject string, data []byte) (*model.Observation, string, error) {
	switch {
	case strings.HasPrefix(subject, "tak.cot."):
		return CoT(data, subject)
	case strings.HasPrefix(subject, "tel.fleet."):
		return Fleet(data, subject)
	case strings.HasPrefix(subject, "evt.world."):
		return World(data, subject)
	case strings.HasPrefix(subject, "intel."):
		return Intel(data, subject)
	case strings.HasPrefix(subject, "det."):
		return Det(data, subject)
	}
	return nil, SkipUnsupported, nil
}
