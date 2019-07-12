package gc_test

import (
	"context"

	"github.com/concourse/concourse/atc/db/dbfakes"
	"github.com/concourse/concourse/atc/gc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckCollector", func() {
	var collector gc.Collector
	var fakeCheckLifecycle *dbfakes.FakeCheckLifecycle

	BeforeEach(func() {
		fakeCheckLifecycle = new(dbfakes.FakeCheckLifecycle)

		collector = gc.NewCheckCollector(fakeCheckLifecycle)
	})

	Describe("Run", func() {
		It("tells the check lifecycle to remove expired checks", func() {
			err := collector.Run(context.TODO())
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeCheckLifecycle.RemoveExpiredChecksCallCount()).To(Equal(1))
		})
	})
})
