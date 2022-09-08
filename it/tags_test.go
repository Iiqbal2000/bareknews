package it

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Iiqbal2000/bareknews"
	v1 "github.com/Iiqbal2000/bareknews/pkg/v1"
	"github.com/Iiqbal2000/bareknews/tags"
	"github.com/gavv/httpexpect/v2"
)

const tagsRoute = "/api/tags"

// TagsTest
type TagsTest struct {
	app http.Handler
}

// Test_Tags is an entry point for the tags endpoint tests.
func Test_Tags(t *testing.T) {
	it := RunDepedencies(t)
	t.Cleanup(it.Teardown)
	shutdown := make(chan os.Signal, 1)

	test := TagsTest{
		app: v1.APIMux(v1.APIMuxConfig{
			Shutdown: shutdown,
			Log:      it.Log,
			DB:       it.DB,
		}),
	}

	t.Run("createTags4xx", test.createTags4xx)
	t.Run("updateTags4xx", test.updateTags4xx)
	t.Run("getTagById4xx", test.getTagById4xx)
	t.Run("deleteTags4xx", test.deleteTags4xx)

	t.Run("crudTags2xx", test.crudTags2xx)
}

func (tt TagsTest) deleteTags4xx(t *testing.T) {
	t.Run("invalid UUID", func(t *testing.T) {
		server := httptest.NewServer(tt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		invalidUUID := "9df76a7f-3d44-4fcf-9222-df65f0774bfb1"

		req := e.DELETE(fmt.Sprintf("%s/%s", tagsRoute, invalidUUID))
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrInvalidUUID.Message())
	})

	t.Run("data is not found", func(t *testing.T) {
		server := httptest.NewServer(tt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		validUUID := "9dc85359-3c16-4e8b-be94-1b2276a1ae4a"

		req := e.DELETE(fmt.Sprintf("%s/%s", tagsRoute, validUUID))
		res := req.Expect()
		res.Status(http.StatusNotFound)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrDataNotFound.Message())
	})
}

func (tt TagsTest) getTagById4xx(t *testing.T) {
	t.Run("invalid UUID", func(t *testing.T) {
		server := httptest.NewServer(tt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		invalidUUID := "9df76a7f-3d44-4fcf-9222-df65f0774bfb1"

		req := e.GET(fmt.Sprintf("%s/%s", tagsRoute, invalidUUID))
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrInvalidUUID.Message())
	})

	t.Run("data is not found", func(t *testing.T) {
		server := httptest.NewServer(tt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		validUUID := "9dc85359-3c16-4e8b-be94-1b2276a1ae4a"

		req := e.GET(fmt.Sprintf("%s/%s", tagsRoute, validUUID))
		res := req.Expect()
		res.Status(http.StatusNotFound)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrDataNotFound.Message())
	})
}

func (tt TagsTest) updateTags4xx(t *testing.T) {
	t.Run("invalid UUID", func(t *testing.T) {
		server := httptest.NewServer(tt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		payload := map[string]interface{}{
			"name": "PHP",
		}

		invalidUUID := "9df76a7f-3d44-4fcf-9222-df65f0774bfb1"

		req := e.PUT(fmt.Sprintf("%s/%s", tagsRoute, invalidUUID)).WithJSON(payload)
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrInvalidUUID.Message())
	})

	t.Run("invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(tt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		// NodeJS tag's UUID
		validUUID := "57cb6822-a459-4d4c-9709-7d0820dc441b"

		req := e.PUT(fmt.Sprintf("%s/%s", tagsRoute, validUUID)).WithJSON("")
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrInvalidJSON.Message())
	})

	t.Run("data is not found", func(t *testing.T) {
		server := httptest.NewServer(tt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		validUUID := "9dc85359-3c16-4e8b-be94-1b2276a1ae4a"
		payload := map[string]interface{}{
			"name": "PHP",
		}

		req := e.PUT(fmt.Sprintf("%s/%s", tagsRoute, validUUID)).WithJSON(payload)
		res := req.Expect()
		res.Status(http.StatusNotFound)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrDataNotFound.Message())
	})

	t.Run("blanking a name field", func(t *testing.T) {
		server := httptest.NewServer(tt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		payload := map[string]interface{}{
			"name": "",
		}

		// NodeJS tag's UUID
		validUUID := "57cb6822-a459-4d4c-9709-7d0820dc441b"

		req := e.PUT(fmt.Sprintf("%s/%s", tagsRoute, validUUID)).WithJSON(payload)
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("fields").Object().ContainsMap(map[string]interface{}{
			"Name": "cannot be blank",
		})
	})
}

func (tt TagsTest) createTags4xx(t *testing.T) {
	t.Run("invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(tt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		req := e.POST(tagsRoute).WithJSON("")
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("error").Equal(bareknews.ErrInvalidJSON.Message())
	})

	t.Run("blanking the name field", func(t *testing.T) {
		server := httptest.NewServer(tt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		payload := map[string]interface{}{
			"name": "",
		}

		req := e.POST(tagsRoute).WithJSON(payload)
		res := req.Expect()
		res.Status(http.StatusBadRequest)

		JSONVal := res.JSON().Object()
		JSONVal.Value("fields").Object().ContainsMap(map[string]interface{}{
			"Name": "cannot be blank",
		})
	})

	t.Run("conflicting the data", func(t *testing.T) {
		server := httptest.NewServer(tt.app)
		defer server.Close()

		e := httpexpect.New(t, server.URL)

		// Go tag is seed data already inside the database.
		payload := map[string]interface{}{
			"name": "Go",
		}

		req := e.POST(tagsRoute).WithJSON(payload)
		res := req.Expect()
		res.Status(http.StatusConflict)

		JSONVal := res.JSON().Object()
		JSONVal.ValueEqual("error", bareknews.ErrDataAlreadyExist.Error())
	})
}

func (tt TagsTest) crudTags2xx(t *testing.T) {
	t.Run("getAllTags200", tt.getAllTags200)
	t.Run("createTag201", tt.createTag201)
	t.Run("updateTag200", tt.updateTag200)
	t.Run("getTagById200", tt.getTagById200)
	t.Run("deleteTag200", tt.deleteTag200)
}

func (tt TagsTest) getAllTags200(t *testing.T) {
	server := httptest.NewServer(tt.app)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	req := e.GET(tagsRoute)
	res := req.Expect()
	res.Status(http.StatusOK)

	JSONVal := res.JSON().Object()
	JSONVal.Value("message").Equal(tags.MsgForGettingAllTags)
	JSONVal.Value("data").Array().Length().Equal(2)
}

func (tt TagsTest) createTag201(t *testing.T) {
	server := httptest.NewServer(tt.app)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	payload := map[string]interface{}{
		"name": "Java",
	}

	req := e.POST(tagsRoute).WithJSON(payload)
	res := req.Expect()
	res.Status(http.StatusCreated)

	JSONVal := res.JSON().Object()
	JSONVal.Value("message").Equal(tags.MsgForCreatedTags)
	JSONVal.Value("data").Object().ContainsKey("id").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("name").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("slug").NotEmpty()
}

func (tt TagsTest) updateTag200(t *testing.T) {
	server := httptest.NewServer(tt.app)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	payload := map[string]interface{}{
		"name": "PHP",
	}

	validId := "9df76a7f-3d44-4fcf-9222-df65f0774bfb"

	req := e.PUT(fmt.Sprintf("%s/%s", tagsRoute, validId)).WithJSON(payload)
	res := req.Expect()
	res.Status(http.StatusOK)

	JSONVal := res.JSON().Object()
	JSONVal.Value("message").Equal(tags.MsgForUpdatingTag)
	JSONVal.Value("data").Object().ContainsKey("id").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("name").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("slug").NotEmpty()
}

func (tt TagsTest) getTagById200(t *testing.T) {
	server := httptest.NewServer(tt.app)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	validId := "9df76a7f-3d44-4fcf-9222-df65f0774bfb"

	req := e.GET(fmt.Sprintf("%s/%s", tagsRoute, validId))
	res := req.Expect()
	res.Status(http.StatusOK)

	JSONVal := res.JSON().Object()
	JSONVal.Value("message").Equal(tags.MsgForGettingTag)
	JSONVal.Value("data").Object().ContainsKey("id").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("name").NotEmpty()
	JSONVal.Value("data").Object().ContainsKey("slug").NotEmpty()
}

func (tt TagsTest) deleteTag200(t *testing.T) {
	server := httptest.NewServer(tt.app)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	validId := "9df76a7f-3d44-4fcf-9222-df65f0774bfb"

	req := e.DELETE(fmt.Sprintf("%s/%s", tagsRoute, validId))
	res := req.Expect()
	res.Status(http.StatusOK)
	JSONVal := res.JSON().Object()
	JSONVal.Value("message").Equal(tags.MsgForDeletingTag)
}
