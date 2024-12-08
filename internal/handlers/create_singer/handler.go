package create_singer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/spanner"
	"github.com/esponges/initial-setup/internal/common"
	"github.com/go-playground/validator/v10"
)

type CreateSingersHandlerImpl struct {
	validator     *validator.Validate
	spannerClient *spanner.Client
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

		_, err = c.spannerClient.ReadWriteTransaction(c.ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			cols := []string{"SingerId", "FirstName", "LastName"}
			err = txn.BufferWrite([]*spanner.Mutation{
				spanner.InsertOrUpdate("Singers", cols, []interface{}{req.SingerId, req.Name, req.LastName}),
			})

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Successfully upserted singer: %s\n", req.Name+" "+req.LastName)
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}
