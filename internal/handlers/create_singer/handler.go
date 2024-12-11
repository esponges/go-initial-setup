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
			log.Println("Upsert failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func (c *CreateSingersHandlerImpl) UpsertSinger(req CreateSingersRequest) (string, error) {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		_, txnErr := c.spannerClient.ReadWriteTransaction(c.ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			cols := []string{"SingerId", "FirstName", "LastName"}
			err := txn.BufferWrite([]*spanner.Mutation{
				spanner.InsertOrUpdate("Singers", cols, []interface{}{req.SingerId, req.Name, req.LastName}),
			})

			if err != nil {
				log.Printf("Attempt %d: BufferWrite failed: %v", attempt+1, err)
				return err // Return the error for retry logic
			}

			return nil // Successful write
		})

		if txnErr == nil {
			return "Successfully upserted singer: " + req.Name + " " + req.LastName, nil
		}

		lastErr = txnErr
		log.Printf("Attempt %d: Transaction failed: %v", attempt+1, txnErr)

		// Optionally implement exponential backoff here before the next attempt
	}

	return "", lastErr // Return the last encountered error after all attempts
}
