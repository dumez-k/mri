package main_test

import (
	"bytes"
	"testing"

	"github.com/paketo-buildpacks/packit"
	main "github.com/paketo-community/mri"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testPlanEntryResolver(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		buffer   *bytes.Buffer
		resolver main.PlanEntryResolver
	)

	it.Before(func() {
		buffer = bytes.NewBuffer(nil)
		resolver = main.NewPlanEntryResolver(main.NewLogEmitter(buffer))
	})

	context("when a buildpack.yml entry is included", func() {
		it("resolves the best plan entry", func() {
			entry := resolver.Resolve([]packit.BuildpackPlanEntry{
				{
					Name:    "mri",
					Version: "other-version",
				},
				{
					Name:    "mri",
					Version: "buildpack-yml-version",
					Metadata: map[string]interface{}{
						"version-source": "buildpack.yml",
					},
				},
			})
			Expect(entry).To(Equal(packit.BuildpackPlanEntry{
				Name:    "mri",
				Version: "buildpack-yml-version",
				Metadata: map[string]interface{}{
					"version-source": "buildpack.yml",
				},
			}))

			Expect(buffer.String()).To(ContainSubstring("    Candidate version sources (in priority order):"))
			Expect(buffer.String()).To(ContainSubstring("      buildpack.yml -> \"buildpack-yml-version\""))
			Expect(buffer.String()).To(ContainSubstring("      <unknown>     -> \"other-version\""))
		})
	})

	context("when entry flags differ", func() {
		context("OR's them together on best plan entry", func() {
			it("has all flags", func() {
				entry := resolver.Resolve([]packit.BuildpackPlanEntry{
					{
						Name:    "mri",
						Version: "buildpack-yml-version",
						Metadata: map[string]interface{}{
							"version-source": "buildpack.yml",
						},
					},
					{
						Name:    "mri",
						Version: "",
						Metadata: map[string]interface{}{
							"build": true,
						},
					},
				})
				Expect(entry).To(Equal(packit.BuildpackPlanEntry{
					Name:    "mri",
					Version: "buildpack-yml-version",
					Metadata: map[string]interface{}{
						"version-source": "buildpack.yml",
						"build":          true,
					},
				}))
			})
		})
	})

	context("when an unknown source entry is included", func() {
		it("resolves the best plan entry", func() {
			entry := resolver.Resolve([]packit.BuildpackPlanEntry{
				{
					Name:    "mri",
					Version: "other-version",
				},
			})
			Expect(entry).To(Equal(packit.BuildpackPlanEntry{
				Name:     "mri",
				Version:  "other-version",
				Metadata: map[string]interface{}{},
			}))
		})
	})
}
