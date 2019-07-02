package api_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/concourse/concourse/atc/api/accessor/accessorfakes"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/db/dbfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Users API", func() {

	var (
		response   *http.Response
		fakeaccess *accessorfakes.FakeAccess
		query      url.Values
	)

	BeforeEach(func() {
		fakeaccess = new(accessorfakes.FakeAccess)
	})

	Context("GET /api/v1/users", func() {

		JustBeforeEach(func() {
			fakeAccessor.CreateReturns(fakeaccess)

			req, err := http.NewRequest("GET", server.URL+"/api/v1/users", nil)
			Expect(err).NotTo(HaveOccurred())

			req.URL.RawQuery = query.Encode()

			response, err = client.Do(req)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authenticated", func() {

			BeforeEach(func() {
				fakeaccess.IsAuthenticatedReturns(true)
			})

			Context("not an admin", func() {

				It("returns 403", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
				})

			})

			Context("being an admin", func() {

				BeforeEach(func() {
					fakeaccess.IsAdminReturns(true)
				})

				It("succeeds", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})

				It("returns Content-Type 'application/json'", func() {
					Expect(response.Header.Get("Content-Type")).To(Equal("application/json"))
				})

				Context("failing to retrieve users", func() {
					BeforeEach(func() {
						dbUserFactory.GetAllUsersReturns(nil, errors.New("no db connection"))
					})

					It("fails", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})

				Context("having no users", func() {
					BeforeEach(func() {
						dbUserFactory.GetAllUsersReturns([]db.User{}, nil)
					})

					It("returns an empty array", func() {
						body, err := ioutil.ReadAll(response.Body)
						Expect(err).NotTo(HaveOccurred())

						Expect(body).To(MatchJSON(`[]`))
					})
				})

				Context("having users", func() {
					BeforeEach(func() {
						user1 := new(dbfakes.FakeUser)
						user1.IDReturns(6)
						user1.NameReturns("bob")
						user1.ConnectorReturns("github")
						user1.LastLoginReturns(time.Unix(10, 0))

						dbUserFactory.GetAllUsersReturns([]db.User{user1}, nil)
					})

					It("returns all users logged in since table creation", func() {
						body, err := ioutil.ReadAll(response.Body)
						Expect(err).NotTo(HaveOccurred())

						Expect(body).To(MatchJSON(`[{
							"id": 6,
							"username": "bob",
							"connector": "github",
							"last_login": "1969-12-31T19:00:10-05:00"
						}]`))
					})

				})

			})

		})

		Context("not authenticated", func() {

			BeforeEach(func() {
				fakeaccess.IsAuthenticatedReturns(false)
			})

			It("returns 401", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})

		})

	})

})