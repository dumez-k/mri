package main

import (
	"os"
	"time"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/cargo"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/postal"
)

func main() {
	buildpackYMLParser := NewBuildpackYMLParser()

	logEmitter := NewLogEmitter(os.Stdout)
	entryResolver := NewPlanEntryResolver(logEmitter)
	dependencyManager := postal.NewService(cargo.NewTransport())
	planRefinery := NewPlanRefinery()
	clock := NewClock(time.Now)
	gem := pexec.NewExecutable("gem")

	packit.Run(
		Detect(buildpackYMLParser),
		Build(entryResolver, dependencyManager, planRefinery, logEmitter, clock, gem),
	)
}
