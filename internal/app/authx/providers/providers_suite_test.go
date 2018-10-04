/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers


import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)


func TestCredentialsPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Providers package suite")
}
