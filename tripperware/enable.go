package tripperware

import (
	"github.com/gol4ng/httpware/v2"
)

// Enable tripperware is used to conditionnaly add a tripperware to a TripperwareStack
// See Skip tripperware to active a tripperware in function of request
func Enable(enable bool, tripperware httpware.Tripperware) httpware.Tripperware {
	if enable {
		return tripperware
	}
	return httpware.NopTripperware
}
