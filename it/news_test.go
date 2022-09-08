package it

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/news"
	v1 "github.com/Iiqbal2000/bareknews/pkg/v1"
	"github.com/gavv/httpexpect/v2"
)

const newsRoute = "/api/news"

// NewsTest
type NewsTest struct {
	app http.Handler
}

// Test_News is an entry point for the news endpoint tests.
func Test_News(t *testing.T) {
	it := RunDepedencies(t)
	t.Cleanup(it.Teardown)
	shutdown := make(chan os.Signal, 1)

	test := NewsTest{
		app: v1.APIMux(v1.APIMuxConfig{
			Shutdown: shutdown,
			Log:      it.Log,
			DB:       it.DB,
		}),
	}

	t.Run("createNews4xx", test.createNews4xx)
	t.Run("getNewsById4xx", test.getNewsById4xx)
	t.Run("updateNews4xx", test.updateNews4xx)
	t.Run("deleteNews4xx", test.deleteNews4xx)

	t.Run("crudNews2xx", test.crudNews2xx)
}

func (nt NewsTest) crudNews2xx(t *testing.T) {
	t.Run("createNews201", nt.createNews201)
	t.Run("getAllNews200", nt.getAllNews200)
	t.Run("getNewsById200", nt.getNewsById200)
	t.Run("updateNews200", nt.updateNews200)
	t.Run("deleteNews200", nt.deleteNews200)
}

func (nt NewsTest) createNews201(t *testing.T) {
	server := httptest.NewServer(nt.app)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	payload := map[string]interface{}{
		"title": "Post 1",
		"body":  "test test test test tes",
		"tags":  []string{"Go", "NodeJS"},
	}

	req := e.POST(newsRoute).WithJSON(payload)
	res := req.Expect()
	res.Status(http.StatusCreated)

	JSONVal := res.JSON().Object()
	JSONVal.Value("message").Equal(news.MsgForCreatedNews)
	JSONVal.Value("data").Object().ContainsKey("id").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("title").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("body").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("status").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("slug").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("tags").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("date_created").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("date_updated").NotEmpty()
}

func (nt NewsTest) createNews4xx(t *testing.T) {
	t.Run("invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(nt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		req := e.POST(newsRoute).WithJSON("")
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrInvalidJSON.Message())
	})

	t.Run("blanking all fields", func(t *testing.T) {
		server := httptest.NewServer(nt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		payload := map[string]interface{}{
			"title": "",
			"body":  "",
		}

		req := e.POST(newsRoute).WithJSON(payload)
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("fields").Object().ContainsMap(map[string]interface{}{
			"Title": "cannot be blank",
			"Body":  "cannot be blank",
		})
	})

	t.Run("conflicting the data", func(t *testing.T) {
		server := httptest.NewServer(nt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		payload := map[string]interface{}{
			"title": "Post 2021",
			"body":  "test test test test tes",
			"tags":  []string{"Go", "NodeJS"},
		}

		req := e.POST(newsRoute).WithJSON(payload)
		res := req.Expect()
		res.Status(http.StatusConflict)

		JSONVal := res.JSON().Object()
		JSONVal.ValueEqual("error", bareknews.ErrDataAlreadyExist.Error())
	})
}

func (nt NewsTest) updateNews200(t *testing.T) {
	server := httptest.NewServer(nt.app)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	payload := map[string]interface{}{
		"title":  "Post 301",
		"body":   "test test test test tes",
		"tags":   []string{"Go", "NodeJS"},
		"status": "publish",
	}

	validUUID := "fdc76bfc-5aa1-4096-a9d3-39719610f987"

	req := e.PUT(fmt.Sprintf("%s/%s", newsRoute, validUUID)).WithJSON(payload)
	res := req.Expect()
	res.Status(http.StatusOK)

	JSONVal := res.JSON().Object()
	JSONVal.Value("message").Equal(news.MsgForUpdatedNews)
	JSONVal.Value("data").Object().ContainsKey("id").NotEmpty()
	JSONVal.Value("data").Object().ValueEqual("title", payload["title"])
	JSONVal.Value("data").Object().ValueEqual("body", payload["body"])
	JSONVal.Value("data").Object().ContainsKey("status").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("slug").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("tags").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("date_created").NotEmpty()
	JSONVal.Value("data").Object().ValueNotEqual("date_updated", JSONVal.Value("data").Object().Value("date_created"))
}

func (nt NewsTest) updateNews4xx(t *testing.T) {
	t.Run("invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(nt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		validUUID := "fdc76bfc-5aa1-4096-a9d3-39719610f987"

		req := e.PUT(fmt.Sprintf("%s/%s", newsRoute, validUUID)).WithJSON("")
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrInvalidJSON.Message())
	})

	t.Run("invalid UUID", func(t *testing.T) {
		server := httptest.NewServer(nt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		invalidUUID := "9df76a7f-3d44-4fcf-9222-df65f0774bfb23"

		payload := map[string]interface{}{
			"title":  "Post 301",
			"body":   "test test test test tes",
			"tags":   []string{"Go", "NodeJS"},
			"status": "publish",
		}

		req := e.PUT(fmt.Sprintf("%s/%s", newsRoute, invalidUUID)).WithJSON(payload)
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrInvalidUUID.Message())
	})

	t.Run("blanking all fields", func(t *testing.T) {
		server := httptest.NewServer(nt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		payload := map[string]interface{}{
			"title": "",
			"body":  "",
		}

		validUUID := "fdc76bfc-5aa1-4096-a9d3-39719610f987"

		req := e.PUT(fmt.Sprintf("%s/%s", newsRoute, validUUID)).WithJSON(payload)
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("fields").Object().ContainsMap(map[string]interface{}{
			"Title": "cannot be blank",
			"Body":  "cannot be blank",
		})
	})

	t.Run("data is not found", func(t *testing.T) {
		server := httptest.NewServer(nt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		payload := map[string]interface{}{
			"title":  "Post 301",
			"body":   "test test test test tes",
			"tags":   []string{"Go", "NodeJS"},
			"status": "publish",
		}

		validUUID := "9dc85359-3c16-4e8b-be94-1b2276a1ae4a"

		req := e.PUT(fmt.Sprintf("%s/%s", newsRoute, validUUID)).WithJSON(payload)
		res := req.Expect()
		res.Status(http.StatusNotFound)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrDataNotFound.Message())
	})
}

func (nt NewsTest) deleteNews200(t *testing.T) {
	server := httptest.NewServer(nt.app)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	validUUID := "fdc76bfc-5aa1-4096-a9d3-39719610f987"

	req := e.DELETE(fmt.Sprintf("%s/%s", newsRoute, validUUID))
	res := req.Expect()
	res.Status(http.StatusOK)
	JSONVal := res.JSON().Object()
	JSONVal.Value("message").Equal(news.MsgForDeletedNews)
}

func (nt NewsTest) deleteNews4xx(t *testing.T) {
	t.Run("invalid UUID", func(t *testing.T) {
		server := httptest.NewServer(nt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		invalidUUID := "9df76a7f-3d44-4fcf-9222-df65f0774bfb23"

		req := e.DELETE(fmt.Sprintf("%s/%s", newsRoute, invalidUUID))
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrInvalidUUID.Message())
	})

	t.Run("data is not found", func(t *testing.T) {
		server := httptest.NewServer(nt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		validUUID := "9dc85359-3c16-4e8b-be94-1b2276a1ae4a"

		req := e.DELETE(fmt.Sprintf("%s/%s", newsRoute, validUUID))
		res := req.Expect()
		res.Status(http.StatusNotFound)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrDataNotFound.Message())
	})
}

func (nt NewsTest) getAllNews200(t *testing.T) {
	server := httptest.NewServer(nt.app)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	req := e.GET(newsRoute)
	res := req.Expect()
	res.Status(http.StatusOK)
	JSONVal := res.JSON().Object()
	JSONVal.Value("message").Equal(news.MsgForGettingAllNews)
	JSONVal.Value("data").Array().Length().Equal(2)
}

func (nt NewsTest) getNewsById200(t *testing.T) {
	server := httptest.NewServer(nt.app)
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	
	validUUID := "fdc76bfc-5aa1-4096-a9d3-39719610f987"
	
	req := e.GET(fmt.Sprintf("%s/%s", newsRoute, validUUID))
	res := req.Expect()
	res.Status(http.StatusOK)
	JSONVal := res.JSON().Object()
	JSONVal.Value("message").Equal(news.MsgForGettingNews)
	JSONVal.Value("data").NotNull()
}

func (nt NewsTest) getNewsById4xx(t *testing.T) {
	t.Run("invalid UUID", func(t *testing.T) {
		server := httptest.NewServer(nt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		invalidUUID := "9df76a7f-3d44-4fcf-9222-df65f0774bfb1"

		req := e.GET(fmt.Sprintf("%s/%s", newsRoute, invalidUUID))
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrInvalidUUID.Message())
	})

	t.Run("data is not found", func(t *testing.T) {
		server := httptest.NewServer(nt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		validUUID := "9dc85359-3c16-4e8b-be94-1b2276a1ae4a"

		req := e.GET(fmt.Sprintf("%s/%s", newsRoute, validUUID))
		res := req.Expect()
		res.Status(http.StatusNotFound)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrDataNotFound.Message())
	})
}