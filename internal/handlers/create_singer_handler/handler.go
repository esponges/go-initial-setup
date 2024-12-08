package create_singer_handler

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
	defer client.Close()

	return &CreateSingersHandlerImpl{
		validator:     validator,
		spannerClient: client,
		ctx:           ctx,
	}
}

func (c *CreateSingersHandlerImpl) CreateSingersHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateSingersRequest
	body, err := common.UnmarshalAndValidateRequest(r, &req, c.validator)
	log.Println("body: " + string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	} else {
		log.Println("Correct Request")

		_, err = c.spannerClient.ReadWriteTransaction(c.ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			// sql := "INSERT Singers (SingerId, FirstName, LastName) VALUES (:singerId, :firstName, :lastName)"

			// none of these work
			// stmt := spanner.Statement{
			// 	SQL: sql,
			// 	Params: map[string]interface{}{
			// 		"singerId":  "13",
			// 		"firstName": "foo",
			// 		"lastName":  "bar",
			// 	},
			// }
			log.Println(string(req.Name), string(req.LastName))

			sql := fmt.Sprintf("INSERT Singers (SingerId, FirstName, LastName) VALUES ('%s', '%s', '%s')",
				req.SingerId, req.Name, req.LastName)

			stmt := spanner.Statement{
				SQL: sql,
			}

			// this sample code works
			// stmt := spanner.Statement{
			// 	SQL: `INSERT Singers (SingerId, FirstName, LastName) VALUES
			// 		(12, 'Melissa', 'Garcia'),
			// 		(13, 'Russell', 'Morales'),
			// 		(14, 'Jacqueline', 'Long'),
			// 		(15, 'Dylan', 'Shaw')`,
			// }

			rowCount, err := txn.Update(ctx, stmt)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%d rows inserted.\n", rowCount)
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}
