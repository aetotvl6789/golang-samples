package storagetransfer

// [START storagetransfer_get_latest_transfer_operation]
import (
	"context"
	"fmt"
	"io"

	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	storagetransferpb "google.golang.org/genproto/googleapis/storagetransfer/v1"
)

func check_latest_transfer_operation(w io.Writer, projectID string, jobName string) (*storagetransferpb.TransferOperation, error) {
	ctx := context.Background()
	client, err := storagetransfer.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storagetransfer.NewClient: %v", err)
	}
	defer client.Close()

	job, err := client.GetTransferJob(ctx, &storagetransferpb.GetTransferJobRequest{
		JobName:   jobName,
		ProjectId: projectID,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to get transfer job: %v", err)
	}

	latestOpName := job.LatestOperationName
	if latestOpName != "" {
		lro, err := client.LROClient.GetOperation(ctx, &longrunningpb.GetOperationRequest{
			Name: latestOpName,
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to get transfer operation: %v", err)
		}
		latestOp := &storagetransferpb.TransferOperation{}
		lro.Metadata.UnmarshalTo(latestOp)

		fmt.Fprintf(w, "The latest transfer operation for job %q is: \n%v", jobName, latestOp)
		return latestOp, nil
	} else {
		fmt.Fprintf(w, "Transfer job %q hasn't run yet, try again later", jobName)
		return nil, nil
	}
}

// [END storagetransfer_get_latest_transfer_operation]
