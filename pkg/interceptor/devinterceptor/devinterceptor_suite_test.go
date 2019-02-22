/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package devinterceptor

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestInterceptorPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Device Interceptor package suite")
}
