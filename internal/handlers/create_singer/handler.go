package create_singer

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/esponges/initial-setup/internal/common"
	"github.com/go-playground/validator/v10"
)

type SpannerClient interface {
	ReadWriteTransaction(context.Context, func(context.Context, *spanner.ReadWriteTransaction) error) (time.Time, error)
}

type SpannerClientWrapper struct {
	client *spanner.Client
}

func (s *SpannerClientWrapper) ReadWriteTransaction(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error) (commitTimestamp time.Time, err error) {
	return s.client.ReadWriteTransaction(ctx, f)
}

type CreateSingersHandlerImpl struct {
	validator     *validator.Validate
	spannerClient SpannerClient
	ctx           context.Context
}

func NewCreateSingersHandler(validator *validator.Validate) *CreateSingersHandlerImpl {
	ctx := context.Background()

	client, err := spanner.NewClient(ctx, os.Getenv("SPANNER_DB"))
	if err != nil {
		log.Fatal(err)
	}

	return &CreateSingersHandlerImpl{
		validator:     validator,
		spannerClient: client,
		ctx:           ctx,
	}
}

func (c *CreateSingersHandlerImpl) CreateSingersHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateSingersRequest
	body, err := common.UnmarshalAndValidateRequest(r, &req, c.validator)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	} else {
		log.Println("Correct Request")

		res, err := c.UpsertSinger(req)
		log.Printf("UpsertSinger returned: res=%s, err=%v", res, err)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func (c *CreateSingersHandlerImpl) UpsertSinger(req CreateSingersRequest) (string, error) {
	_, err := c.spannerClient.ReadWriteTransaction(c.ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		cols := []string{"SingerId", "FirstName", "LastName"}
		err := txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Singers", cols, []interface{}{req.SingerId, req.Name, req.LastName}),
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return "Successfully upserted singer: " + req.Name + " " + req.LastName, nil
}
