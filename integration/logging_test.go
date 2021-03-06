package integration

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testLogging(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when the buildpack is run with pack build", func() {
		var (
			image occam.Image
			name  string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
		})

		it("logs useful information for the user", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithNoPull().
				WithBuildpacks(mriBuildpack, buildPlanBuildpack).
				Execute(name, filepath.Join("testdata", "simple_app"))
			Expect(err).ToNot(HaveOccurred(), logs.String)

			buildpackVersion, err := GetGitVersion()
			Expect(err).ToNot(HaveOccurred())

			Expect(logs).To(ContainLines(
				fmt.Sprintf("MRI Buildpack %s", buildpackVersion),
				"  Resolving MRI version",
				"    Candidate version sources (in priority order):",
				"      buildpack.yml -> \"2.7.x\"",
				"      <unknown>     -> \"*\"",
				"",
				MatchRegexp(`    Selected MRI version \(using buildpack\.yml\): 2\.7\.\d+`),
				"",
				"  Executing build process",
				MatchRegexp(`    Installing MRI 2\.7\.\d+`),
				MatchRegexp(`      Completed in \d+\.\d+`),
				"",
				"  Configuring environment",
				MatchRegexp(`    GEM_PATH -> "/home/vcap/.gem/ruby/2\.7\.\d+:/layers/paketo-community_mri/mri/lib/ruby/gems/2\.7\.\d+"`),
			))
		})
	})
}
