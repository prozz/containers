package containers_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"git.naspersclassifieds.com/olxeu/monetization/jobs/new-jobs-app/backend/platform/containers"
)

func TestDynamoContainer(t *testing.T) {
	t.SkipNow()

	container, terminateFn := containers.Start(t, context.Background(),
		"amazon/dynamodb-local", []string{"8000/tcp"}, wait.ForListeningPort("8000/tcp"))
	defer terminateFn()

	endpointURL, err := container.Endpoint(context.Background(), "http")
	if err != nil {
		t.Fatal(err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("local"),
		Endpoint: aws.String(endpointURL),
	})
	if err != nil {
		t.Fatal(err)
	}

	db := dynamodb.New(sess)

	for i := 0; i < 10; i++ {
		start := time.Now()
		output, err := db.ListTables(&dynamodb.ListTablesInput{})
		assert.NoError(t, err)
		assert.NotEmpty(t, output)
		t.Logf("list tables: %s", time.Since(start))
	}
}

