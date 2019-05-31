/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package interceptor

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Permissions", func() {
	ginkgo.Context("empty", func() {
		p := Permission{}

		ginkgo.It("allows empty primitives", func() {
			valid := p.Valid([]string{})
			gomega.Expect(valid).To(gomega.BeTrue())
		})

		ginkgo.It("allows any primitive", func() {
			valid := p.Valid([]string{"anyone"})
			gomega.Expect(valid).To(gomega.BeTrue())
		})

		ginkgo.It("allows all primitives", func() {
			valid := p.Valid([]string{"some", "primitive"})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
	})

	ginkgo.Context("with a must primitives", func() {
		validPrimitive := "validPrimitive"
		p := Permission{Must: []string{validPrimitive}}

		ginkgo.It("doesn't allow empty primitives", func() {
			valid := p.Valid([]string{})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("doesn't allow invalid primitive", func() {
			valid := p.Valid([]string{"anyone"})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("doesn't allow all invalid primitives", func() {
			valid := p.Valid([]string{"some", "primitive"})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("allow with a valid primitive", func() {
			valid := p.Valid([]string{validPrimitive})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
		ginkgo.It("allow with a valid primitive and another invalid", func() {
			valid := p.Valid([]string{validPrimitive, "anyone"})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
	})

	ginkgo.Context("with multiple must primitives", func() {
		validPrimitive := "validPrimitive"
		validPrimitive1 := "anotherValidPrimitive"
		p := Permission{Must: []string{validPrimitive, validPrimitive1}}

		ginkgo.It("doesn't allow empty primitives", func() {
			valid := p.Valid([]string{})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("doesn't allow invalid primitive", func() {
			valid := p.Valid([]string{"anyone"})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("doesn't allow all invalid primitives", func() {
			valid := p.Valid([]string{"some", "primitive"})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("doesn't allow with a valid primitive", func() {
			valid := p.Valid([]string{validPrimitive})
			gomega.Expect(valid).To(gomega.BeFalse())
		})
		ginkgo.It("doesn't allow with a valid primitive and another invalid", func() {
			valid := p.Valid([]string{validPrimitive, "anyone"})
			gomega.Expect(valid).To(gomega.BeFalse())
		})
		ginkgo.It("allow with two valid primitives", func() {
			valid := p.Valid([]string{validPrimitive, validPrimitive1})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
		ginkgo.It("allow with two valid primitives and another invalid", func() {
			valid := p.Valid([]string{validPrimitive, validPrimitive1, "anyone"})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
	})

	ginkgo.Context("with a should primitive", func() {
		validPrimitive := "validPrimitive"
		p := Permission{Should: []string{validPrimitive}}

		ginkgo.It("doesn't allow empty primitives", func() {
			valid := p.Valid([]string{})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("doesn't allow invalid primitive", func() {
			valid := p.Valid([]string{"anyone"})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("doesn't allow all invalid primitives", func() {
			valid := p.Valid([]string{"some", "primitive"})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("allow with a valid primitive", func() {
			valid := p.Valid([]string{validPrimitive})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
		ginkgo.It("allow with a valid primitive and another invalid", func() {
			valid := p.Valid([]string{validPrimitive, "anyone"})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
	})

	ginkgo.Context("with multiple should primitives", func() {
		validPrimitive := "validPrimitive"
		validPrimitive1 := "anotherValidPrimitive"
		p := Permission{Should: []string{validPrimitive, validPrimitive1}}

		ginkgo.It("doesn't allow empty primitives", func() {
			valid := p.Valid([]string{})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("doesn't allow invalid primitive", func() {
			valid := p.Valid([]string{"anyone"})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("doesn't allow all invalid primitives", func() {
			valid := p.Valid([]string{"some", "primitive"})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("allow with a valid primitive", func() {
			valid := p.Valid([]string{validPrimitive})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
		ginkgo.It("allow with a valid primitive and another invalid", func() {
			valid := p.Valid([]string{validPrimitive, "anyone"})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
		ginkgo.It("allow with two valid primitives", func() {
			valid := p.Valid([]string{validPrimitive, validPrimitive1})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
		ginkgo.It("allow with two valid primitives and another invalid", func() {
			valid := p.Valid([]string{validPrimitive, validPrimitive1, "anyone"})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
	})

	ginkgo.Context("with a mustNot primitive", func() {
		invalidPrimitive := "invalidPrimitive"
		p := Permission{MustNot: []string{invalidPrimitive}}

		ginkgo.It("allow empty primitives", func() {
			valid := p.Valid([]string{})
			gomega.Expect(valid).To(gomega.BeTrue())
		})

		ginkgo.It("allow any primitive", func() {
			valid := p.Valid([]string{"anyone"})
			gomega.Expect(valid).To(gomega.BeTrue())
		})

		ginkgo.It("allow all primitives", func() {
			valid := p.Valid([]string{"some", "primitive"})
			gomega.Expect(valid).To(gomega.BeTrue())
		})

		ginkgo.It("doesn't allow with a invalid primitive", func() {
			valid := p.Valid([]string{invalidPrimitive})
			gomega.Expect(valid).To(gomega.BeFalse())
		})
		ginkgo.It("doesn't allow with a invalid primitive and another invalid", func() {
			valid := p.Valid([]string{"anyone", invalidPrimitive})
			gomega.Expect(valid).To(gomega.BeFalse())
		})
	})

	ginkgo.Context("with multiple mustNot primitives", func() {
		invalidPrimitive := "invalidPrimitive"
		invalidPrimitive1 := "anotherValidPrimitive"
		p := Permission{MustNot: []string{invalidPrimitive, invalidPrimitive1}}

		ginkgo.It("allow empty primitives", func() {
			valid := p.Valid([]string{})
			gomega.Expect(valid).To(gomega.BeTrue())
		})

		ginkgo.It("allow any primitive", func() {
			valid := p.Valid([]string{"anyone"})
			gomega.Expect(valid).To(gomega.BeTrue())
		})

		ginkgo.It("allow all primitives", func() {
			valid := p.Valid([]string{"some", "primitive"})
			gomega.Expect(valid).To(gomega.BeTrue())
		})

		ginkgo.It("doesn't allow with a invalid primitive", func() {
			valid := p.Valid([]string{invalidPrimitive})
			gomega.Expect(valid).To(gomega.BeFalse())
		})
		ginkgo.It("doesn't allow with a invalid primitive and another invalid", func() {
			valid := p.Valid([]string{invalidPrimitive, "anyone"})
			gomega.Expect(valid).To(gomega.BeFalse())
		})
		ginkgo.It("doesn't allow with two invalid primitives", func() {
			valid := p.Valid([]string{invalidPrimitive, invalidPrimitive1})
			gomega.Expect(valid).To(gomega.BeFalse())
		})
		ginkgo.It("doesn't allow with two invalid primitives and another primitive", func() {
			valid := p.Valid([]string{invalidPrimitive, invalidPrimitive1, "anyone"})
			gomega.Expect(valid).To(gomega.BeFalse())
		})
	})

	ginkgo.Context("with multiple primitives", func() {
		mustPrimitive := "mustPrimitive"
		mustPrimitive1 := "anotherMustPrimitive"
		shouldPrimitive := "shouldPrimitive"
		shouldPrimitive1 := "anotherShouldPrimitive"
		mustNotPrimitive := "invalidPrimitive"
		mustNotPrimitive1 := "anotherValidPrimitive"

		p := Permission{
			Must:    []string{mustPrimitive, mustPrimitive1},
			Should:  []string{shouldPrimitive, shouldPrimitive1},
			MustNot: []string{mustNotPrimitive, mustNotPrimitive1},
		}

		ginkgo.It("doesn't allow empty primitives", func() {
			valid := p.Valid([]string{})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

		ginkgo.It("doesn't allow only must primitives", func() {
			valid := p.Valid([]string{mustPrimitive, mustPrimitive1})
			gomega.Expect(valid).To(gomega.BeFalse())
		})
		ginkgo.It("doesn't allow only a should primitive", func() {
			valid := p.Valid([]string{shouldPrimitive, shouldPrimitive1})
			gomega.Expect(valid).To(gomega.BeFalse())
		})
		ginkgo.It("allow with must primitives and a should primitive", func() {
			valid := p.Valid([]string{shouldPrimitive, mustPrimitive, mustPrimitive1})
			gomega.Expect(valid).To(gomega.BeTrue())
		})
		ginkgo.It("doesn't allow with mustNot primitives", func() {
			valid := p.Valid([]string{shouldPrimitive, mustPrimitive, mustPrimitive1, mustNotPrimitive1})
			gomega.Expect(valid).To(gomega.BeFalse())
		})

	})
})
