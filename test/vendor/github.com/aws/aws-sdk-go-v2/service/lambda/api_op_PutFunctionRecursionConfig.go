// Code generated by smithy-go-codegen DO NOT EDIT.

package lambda

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Sets your function's [recursive loop detection] configuration.
//
// When you configure a Lambda function to output to the same service or resource
// that invokes the function, it's possible to create an infinite recursive loop.
// For example, a Lambda function might write a message to an Amazon Simple Queue
// Service (Amazon SQS) queue, which then invokes the same function. This
// invocation causes the function to write another message to the queue, which in
// turn invokes the function again.
//
// Lambda can detect certain types of recursive loops shortly after they occur.
// When Lambda detects a recursive loop and your function's recursive loop
// detection configuration is set to Terminate , it stops your function being
// invoked and notifies you.
//
// [recursive loop detection]: https://docs.aws.amazon.com/lambda/latest/dg/invocation-recursion.html
func (c *Client) PutFunctionRecursionConfig(ctx context.Context, params *PutFunctionRecursionConfigInput, optFns ...func(*Options)) (*PutFunctionRecursionConfigOutput, error) {
	if params == nil {
		params = &PutFunctionRecursionConfigInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "PutFunctionRecursionConfig", params, optFns, c.addOperationPutFunctionRecursionConfigMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*PutFunctionRecursionConfigOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type PutFunctionRecursionConfigInput struct {

	// The name or ARN of the Lambda function.
	//
	// Name formats
	//
	//   - Function name – my-function .
	//
	//   - Function ARN – arn:aws:lambda:us-west-2:123456789012:function:my-function .
	//
	//   - Partial ARN – 123456789012:function:my-function .
	//
	// The length constraint applies only to the full ARN. If you specify only the
	// function name, it is limited to 64 characters in length.
	//
	// This member is required.
	FunctionName *string

	// If you set your function's recursive loop detection configuration to Allow ,
	// Lambda doesn't take any action when it detects your function being invoked as
	// part of a recursive loop. We recommend that you only use this setting if your
	// design intentionally uses a Lambda function to write data back to the same
	// Amazon Web Services resource that invokes it.
	//
	// If you set your function's recursive loop detection configuration to Terminate ,
	// Lambda stops your function being invoked and notifies you when it detects your
	// function being invoked as part of a recursive loop.
	//
	// By default, Lambda sets your function's configuration to Terminate .
	//
	// If your design intentionally uses a Lambda function to write data back to the
	// same Amazon Web Services resource that invokes the function, then use caution
	// and implement suitable guard rails to prevent unexpected charges being billed to
	// your Amazon Web Services account. To learn more about best practices for using
	// recursive invocation patterns, see [Recursive patterns that cause run-away Lambda functions]in Serverless Land.
	//
	// [Recursive patterns that cause run-away Lambda functions]: https://serverlessland.com/content/service/lambda/guides/aws-lambda-operator-guide/recursive-runaway
	//
	// This member is required.
	RecursiveLoop types.RecursiveLoop

	noSmithyDocumentSerde
}

type PutFunctionRecursionConfigOutput struct {

	// The status of your function's recursive loop detection configuration.
	//
	// When this value is set to Allow and Lambda detects your function being invoked
	// as part of a recursive loop, it doesn't take any action.
	//
	// When this value is set to Terminate and Lambda detects your function being
	// invoked as part of a recursive loop, it stops your function being invoked and
	// notifies you.
	RecursiveLoop types.RecursiveLoop

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationPutFunctionRecursionConfigMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestjson1_serializeOpPutFunctionRecursionConfig{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpPutFunctionRecursionConfig{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "PutFunctionRecursionConfig"); err != nil {
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
	if err = addOpPutFunctionRecursionConfigValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opPutFunctionRecursionConfig(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opPutFunctionRecursionConfig(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "PutFunctionRecursionConfig",
	}
}