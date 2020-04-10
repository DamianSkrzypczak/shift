package autoapi

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/DamianSkrzypczak/shift"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

type successHandler func(Operation, shift.RequestContext, interface{}) error
type businessErrorHandler func(error, Operation, shift.RequestContext, interface{}) error
type validationErrorHandler func(error, *gojsonschema.Result, shift.RequestContext) error
type internalErrorHandler func(error, shift.RequestContext)

type handlers struct {
	successHandler         successHandler
	businessErrorHandler   businessErrorHandler
	validationErrorHandler validationErrorHandler
	internalErrorHandler   internalErrorHandler
}

type Deserializer func(v interface{}) error

type DeserializerFactory func(rc shift.RequestContext) Deserializer

type ResourceAPI struct {
	Domain              *shift.Domain
	resourceURL         string
	handlers            handlers
	deserializerFactory DeserializerFactory
}

func NewResourceAPI(d *shift.Domain) *ResourceAPI {
	return &ResourceAPI{
		Domain:      d,
		resourceURL: "/{resourceID}",
		handlers: handlers{
			successHandler:         newDefaultSuccessHandler(d),
			businessErrorHandler:   newDefaultBusinessErrorHandler(d),
			validationErrorHandler: newDefaultValidationErrorHandler(d),
			internalErrorHandler:   newDefaultInternalErrorHandler(d),
		},
		deserializerFactory: newDefaultDeserializerFactory(d),
	}
}

func (api *ResourceAPI) ResourceURL(id string) string {
	if id != "" {
		return path.Join(api.Domain.Path, id)
	}

	return api.Domain.Path
}

func (api *ResourceAPI) SuccessHandler(handler successHandler) {
	api.handlers.successHandler = handler
}
func (api *ResourceAPI) BusinessErrorHandler(handler businessErrorHandler) {
	api.handlers.businessErrorHandler = handler
}
func (api *ResourceAPI) ValidationErrorHandler(handler validationErrorHandler) {
	api.handlers.validationErrorHandler = handler
}
func (api *ResourceAPI) InternalErrorHandler(handler internalErrorHandler) {
	api.handlers.internalErrorHandler = handler
}
func (api *ResourceAPI) DeserializerFactory(factory DeserializerFactory) {
	api.deserializerFactory = factory
}

func (api *ResourceAPI) newJSONSchemaValidator(loader gojsonschema.JSONLoader) shift.Middleware {
	return func(next shift.Handler) shift.Handler {
		return func(rc shift.RequestContext) {
			body, err := rc.Request.BodyCopy()

			if err != nil {
				api.handlers.internalErrorHandler(err, rc)
				return
			}

			validationResult, err := gojsonschema.Validate(
				loader,
				gojsonschema.NewBytesLoader(body),
			)

			if err != nil || !validationResult.Valid() {
				if err = api.handlers.validationErrorHandler(err, validationResult, rc); err != nil {
					api.handlers.internalErrorHandler(err, rc)
					return
				}

				return
			}

			next(rc)
		}
	}
}

func (api *ResourceAPI) List(
	dataProvider func(params shift.QueryParameters) (interface{}, error),
) {
	api.Domain.Router.
		Get("/", func(rc shift.RequestContext) {
			v, err := dataProvider(rc.Request.QueryParameters)
			api.runSubHandlers(List, rc, v, err)
		})
}
func (api *ResourceAPI) Create(
	loader gojsonschema.JSONLoader,
	dataReceiver func(deserializer Deserializer, params shift.QueryParameters) (interface{}, error),
) {
	api.Domain.Router.
		With(api.newJSONSchemaValidator(loader)).
		Post("/", func(rc shift.RequestContext) {
			v, err := dataReceiver(api.deserializerFactory(rc), rc.Request.QueryParameters)
			api.runSubHandlers(Create, rc, v, err)
		})
}
func (api *ResourceAPI) Read(
	dataProvider func(id string, params shift.QueryParameters) (interface{}, error),
) {
	api.Domain.Router.
		Get(api.resourceURL, func(rc shift.RequestContext) {
			v, err := dataProvider(rc.Request.URLParam("resourceID"), rc.Request.QueryParameters)
			api.runSubHandlers(Read, rc, v, err)
		})
}

func (api *ResourceAPI) Replace(
	loader gojsonschema.JSONLoader,
	dataReceiver func(deserializer Deserializer, id string, params shift.QueryParameters) (interface{}, error),
) {
	api.Domain.Router.
		With(api.newJSONSchemaValidator(loader)).
		Put(api.resourceURL, func(rc shift.RequestContext) {
			v, err := dataReceiver(api.deserializerFactory(rc), rc.Request.URLParam("resourceID"), rc.Request.QueryParameters)
			api.runSubHandlers(Replace, rc, v, err)
		})
}
func (api *ResourceAPI) Update(
	loader gojsonschema.JSONLoader,
	dataReceiver func(deserializer Deserializer, id string, params shift.QueryParameters) (interface{}, error),
) {
	api.Domain.Router.
		With(api.newJSONSchemaValidator(loader)).
		Patch(api.resourceURL, func(rc shift.RequestContext) {
			v, err := dataReceiver(api.deserializerFactory(rc), rc.Request.URLParam("resourceID"), rc.Request.QueryParameters)
			api.runSubHandlers(Update, rc, v, err)
		})
}
func (api *ResourceAPI) Delete(
	dataReceiver func(id string, params shift.QueryParameters) (interface{}, error),
) {
	api.Domain.Router.
		Delete(api.resourceURL, func(rc shift.RequestContext) {
			v, err := dataReceiver(rc.Request.URLParam("resourceID"), rc.Request.QueryParameters)
			api.runSubHandlers(Delete, rc, v, err)
		})
}

func (api *ResourceAPI) runSubHandlers(op Operation, rc shift.RequestContext, v interface{}, err error) {
	if err != nil {
		if err = api.handlers.businessErrorHandler(err, op, rc, v); err != nil {
			api.handlers.internalErrorHandler(err, rc)
			return
		}

		return
	}

	if err = api.handlers.successHandler(op, rc, v); err != nil {
		api.handlers.internalErrorHandler(err, rc)
	}
}

type Operation string

const (
	List    Operation = "list"
	Create  Operation = "create"
	Read    Operation = "read"
	Replace Operation = "replace"
	Update  Operation = "update"
	Delete  Operation = "delete"
)

func newDefaultSuccessHandler(d *shift.Domain) successHandler {
	return func(op Operation, rc shift.RequestContext, v interface{}) error {
		switch op {
		case List:
			return rc.Response.WithJSON(v, http.StatusOK)
		case Create:
			return rc.Response.WithJSON(v, http.StatusCreated)
		case Read:
			return rc.Response.WithJSON(v, http.StatusOK)
		case Replace:
			rc.Response.SetStatusCode(http.StatusNoContent)
		case Update:
			rc.Response.SetStatusCode(http.StatusNoContent)
		case Delete:
			rc.Response.SetStatusCode(http.StatusNoContent)
		}

		return nil
	}
}
func newDefaultBusinessErrorHandler(d *shift.Domain) businessErrorHandler {
	return func(err error, op Operation, rc shift.RequestContext, v interface{}) error {
		return err // promote error to InternalErrorHandler.
	}
}
func newDefaultValidationErrorHandler(d *shift.Domain) validationErrorHandler {
	type ValidationError struct {
		Context     string `json:"context"`
		Description string `json:"description"`
		LogID       string `json:"LogID,omitempty"`
	}

	type ValidationErrors struct {
		Errors []ValidationError `json:"errors"`
	}

	return func(err error, result *gojsonschema.Result, rc shift.RequestContext) error {
		payload := ValidationErrors{}

		if err != nil {
			LogID := ksuid.New().String()
			payload.Errors = append(
				payload.Errors,
				ValidationError{
					"",
					"Abnormal request body.",
					LogID,
				},
			)

			logrus.WithField("LogID", LogID).Error(err)
		} else {
			for _, desc := range result.Errors() {
				payload.Errors = append(
					payload.Errors,
					ValidationError{
						desc.Context().String(),
						desc.Description(),
						"",
					},
				)
			}
		}

		if err := rc.Response.WithJSON(payload, http.StatusBadRequest); err != nil {
			return err
		}

		return nil
	}
}
func newDefaultInternalErrorHandler(d *shift.Domain) internalErrorHandler {
	type InternalErrorResponse struct {
		LogID string `json:"LogID"`
	}

	return func(err error, rc shift.RequestContext) {
		LogID := ksuid.New().String()

		if responseErr := rc.Response.WithJSON(InternalErrorResponse{LogID}, http.StatusInternalServerError); responseErr != nil {
			panic(responseErr)
		}

		logrus.WithField("LogID", LogID).Error(err)
	}
}
func newDefaultDeserializerFactory(d *shift.Domain) DeserializerFactory {
	return func(rc shift.RequestContext) Deserializer {
		return func(v interface{}) error {
			body, err := rc.Request.BodyCopy()
			if err != nil {
				return err
			}

			return json.Unmarshal(body, v)
		}
	}
}
