package transitive

import "example.com/facade"

func Build(meter facade.Meter) {
	_, _ = meter.Int64Counter("app.raw") // want `instrument must be created by a package listed in allowed-packages`
	create := meter.Int64Counter         // want `instrument must be created by a package listed in allowed-packages`
	_, _ = create("app.raw")
}
