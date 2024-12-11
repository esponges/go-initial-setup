package create_singer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/go-playground/validator/v10"
)

// Mock Spanner Client to fully implement the spanner.Client interface
type mockSpannerClient struct {
	readWriteTransactionFunc func(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error) (time.Time, error)
}

// Implement methods required by spanner.Client interface
func (m *mockSpannerClient) ReadWriteTransaction(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error) (time.Time, error) {
	if m.readWriteTransactionFunc != nil {
		return m.readWriteTransactionFunc(ctx, f)
	}
	return time.Time{}, nil
}

// // Additional required methods to satisfy the interface (minimal implementations)
// func (m *mockSpannerClient) Close() {}
// func (m *mockSpannerClient) Single() *spanner.ReadOnlyTransaction {
// 	return nil
// }
// func (m *mockSpannerClient) ReadOnlyTransaction() *spanner.ReadOnlyTransaction {
// 	return nil
// }
// func (m *mockSpannerClient) BatchReadOnlyTransaction(ctx context.Context, tb spanner.TimestampBound) (*spanner.BatchReadOnlyTransaction, error) {
// 	return nil, nil
// }
// func (m *mockSpannerClient) ReadOnlyTransactionWithTimestamp(ts time.Time) *spanner.ReadOnlyTransaction {
// 	return nil
// }
// func (m *mockSpannerClient) Apply(ctx context.Context, ms []*spanner.Mutation) (time.Time, error) {
// 	return time.Time{}, nil
// }
// func (m *mockSpannerClient) ReadWriteTransactionWithOptions(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error, opts spanner.TransactionOptions) (time.Time, error) {
// 	return time.Time{}, nil
// }

// // Mock ReadWriteTransaction to simulate different scenarios
// type mockReadWriteTransaction struct {
// 	bufferWriteFunc func(mutations []*spanner.Mutation) error
// }

// func (m *mockReadWriteTransaction) BufferWrite(mutations []*spanner.Mutation) error {
// 	if m.bufferWriteFunc != nil {
// 		return m.bufferWriteFunc(mutations)
// 	}
// 	return nil
// }

// Test cases for CreateSingersHandler
func TestCreateSingersHandler(t *testing.T) {
	testCases := []struct {
		name               string
		requestBody        CreateSingersRequest
		mockTransactionFn  func(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error) (time.Time, error)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "Successful Singer Creation",
			requestBody: CreateSingersRequest{
				SingerId: "123",
				Name:     "John",
				LastName: "Doe",
			},
			mockTransactionFn: func(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error) (time.Time, error) {
				// Simulate successful transaction
				err := f(ctx, &spanner.ReadWriteTransaction{})
				return time.Now(), err
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"singer_id":"123","name":"John","last_name":"Doe"}`,
		},
		{
			name: "Transaction Failure",
			requestBody: CreateSingersRequest{
				SingerId: "456",
				Name:     "Jane",
				LastName: "Smith",
			},
			mockTransactionFn: func(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error) (time.Time, error) {
				// Simulate transaction failure
				return time.Time{}, errors.New("transaction failed")
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "transaction failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock validator
			validate := validator.New()

			// Create mock Spanner client
			ctx := context.Background()
			mockClient := &mockSpannerClient{
				readWriteTransactionFunc: func(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error) (time.Time, error) {
					return tc.mockTransactionFn(ctx, f)
				},
			}

			// Create handler with mock dependencies
			handler := &CreateSingersHandlerImpl{
				validator:     validate,
				spannerClient: mockClient,
				ctx:           ctx,
			}

			// Prepare request body
			jsonBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Create HTTP request
			req, err := http.NewRequest("POST", "/create-singer", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call the handler
			handler.CreateSingersHandler(w, req)

			// Check response status code
			if w.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, w.Code)
			}

			// Check response body
			// todo: fix expected bodies
			responseBody := w.Body.String()
			if tc.expectedStatusCode == http.StatusOK {
				// For successful case, check the exact JSON
				fmt.Println(responseBody)
				if responseBody != tc.expectedBody {
					t.Errorf("Expected body %s, got %s", tc.expectedBody, responseBody)
				}
			} else {
				// For error cases, check error message
				if responseBody != tc.expectedBody {
					t.Errorf("Expected error message %s, got %s", tc.expectedBody, responseBody)
				}
			}
		})
	}
}

// Test UpsertSinger method separately
func TestUpsertSinger(t *testing.T) {
	testCases := []struct {
		name              string
		requestBody       CreateSingersRequest
		mockTransactionFn func(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error) (time.Time, error)
		expectedResult    string
		expectedError     bool
	}{
		{
			name: "Successful Upsert",
			requestBody: CreateSingersRequest{
				SingerId: "123",
				Name:     "John",
				LastName: "Doe",
			},
			mockTransactionFn: func(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error) (time.Time, error) {
				// Simulate successful transaction
				err := f(ctx, &spanner.ReadWriteTransaction{})
				return time.Now(), err
			},
			expectedResult: "Successfully upserted singer: John Doe",
			expectedError:  false,
		},
		{
			name: "Upsert Failure",
			requestBody: CreateSingersRequest{
				SingerId: "456",
				Name:     "Jane",
				LastName: "Smith",
			},
			mockTransactionFn: func(ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error) (time.Time, error) {
				// Simulate transaction failure
				return time.Time{}, errors.New("upsert failed")
			},
			expectedResult: "",
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock Spanner client
			mockClient := &mockSpannerClient{
				readWriteTransactionFunc: tc.mockTransactionFn,
			}

			// Create handler with mock dependencies
			handler := &CreateSingersHandlerImpl{
				spannerClient: mockClient,
				ctx:           context.Background(),
			}

			// Call UpsertSinger
			result, err := handler.UpsertSinger(tc.requestBody)

			// Check result
			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected an error, got nil")
				}
				if result != "" {
					t.Errorf("Expected empty result, got %s", result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tc.expectedResult {
					t.Errorf("Expected result %s, got %s", tc.expectedResult, result)
				}
			}
		})
	}
}
