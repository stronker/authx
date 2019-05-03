/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package inventory

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Inventory provider", func(){

	sp := NewMockupInventoryProvider()
	RunTest(sp)

})