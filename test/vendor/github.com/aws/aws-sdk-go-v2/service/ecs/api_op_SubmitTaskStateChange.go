// Code generated by smithy-go-codegen DO NOT EDIT.

package ecs

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"time"
)

// This action is only used by the Amazon ECS agent, and it is not intended for
// use outside of the agent.
//
// Sent to acknowledge that a task changed states.
func (c *Client) SubmitTaskStateChange(ctx context.Context, params *SubmitTaskStateChangeInput, optFns ...func(*Options)) (*SubmitTaskStateChangeOutput, error) {
	if params == nil {
		params = &SubmitTaskStateChangeInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "SubmitTaskStateChange", params, optFns, c.addOperationSubmitTaskStateChangeMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*SubmitTaskStateChangeOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type SubmitTaskStateChangeInput struct {

	// Any attachments associated with the state change request.
	Attachments []types.AttachmentStateChange

	// The short name or full Amazon Resource Name (ARN) of the cluster that hosts the
	// task.
	Cluster *string

	// Any containers that's associated with the state change request.
	Containers []types.ContainerStateChange

	// The Unix timestamp for the time when the task execution stopped.
	ExecutionStoppedAt *time.Time

	// The details for the managed agent that's associated with the task.
	ManagedAgents []types.ManagedAgentStateChange

	// The Unix timestamp for the time when the container image pull started.
	PullStartedAt *time.Time

	// The Unix timestamp for the time when the container image pull completed.
	PullStoppedAt *time.Time

	// The reason for the state change request.
	Reason *string

	// The status of the state change request.
	Status *string

	// The task ID or full ARN of the task in the state change request.
	Task *string

	noSmithyDocumentSerde
}

type SubmitTaskStateChangeOutput struct {

	// Acknowledgement of the state change.
	Acknowledgment *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationSubmitTaskStateChangeMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpSubmitTaskStateChange{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpSubmitTaskStateChange{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "SubmitTaskStateChange"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpSubmitTaskStateChangeValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opSubmitTaskStateChange(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opSubmitTaskStateChange(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "SubmitTaskStateChange",
	}
}