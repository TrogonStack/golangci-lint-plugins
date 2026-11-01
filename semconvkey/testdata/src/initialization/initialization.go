package initialization

import "example.com/appsemconv"

func Build() {
	var value appsemconv.TierAttr
	value = appsemconv.NewTierAttr("standard")
	_ = value.KeyValue()
}
