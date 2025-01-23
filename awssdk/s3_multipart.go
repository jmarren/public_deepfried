package awssdk

import (
	// "bytes"
	"context"
	"fmt"

	// "errors"
	// "fmt"
	"io"
	// "os"
	// "time"
	//
	"github.com/aws/aws-sdk-go-v2/aws"
	// "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	// "github.com/aws/aws-sdk-go-v2/service/s3/types"
	// "github.com/jmarren/deepfried/awssdk"
	"github.com/jmarren/deepfried/util"
	// "github.com/jmarren/deepfried/util"
)

func NewMultipart(ctx context.Context, dest string, contentType string) *string {
	fmt.Println("NEW MULTIPART")
	fmt.Printf("dest: %s\n", dest)

	upload, err := s3Client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:      loopS3PublicBucketName,
		Key:         aws.String(dest),
		ContentType: &contentType,
	})
	util.EMsg(err, "creating s3 multipart upload")

	uploadId := upload.UploadId
	fmt.Printf("AWS UPLOAD ID: %s\n", *uploadId)
	return uploadId
}

func UploadPart(ctx context.Context, dest string, uploadId *string, partNumber *int32, dataReader io.Reader) (types.CompletedPart, error) {

	params := &s3.UploadPartInput{
		Bucket:     loopS3PublicBucketName,
		Key:        aws.String(dest),
		UploadId:   uploadId,
		PartNumber: partNumber,
		Body:       dataReader,
	}
	output, err := s3Client.UploadPart(ctx, params)

	part := types.CompletedPart{
		PartNumber: partNumber,
		ETag:       output.ETag,
	}

	fmt.Printf("upload part output: %v\n", output)
	return part, err
}

func CompleteUpload(ctx context.Context, dest string, uploadId *string, completedUpload *types.CompletedMultipartUpload) error {

	fmt.Printf("dest: %s\n", dest)

	params := &s3.CompleteMultipartUploadInput{
		Bucket:          loopS3PublicBucketName,
		Key:             aws.String(dest),
		UploadId:        uploadId,
		MultipartUpload: completedUpload,
	}

	_, err := s3Client.CompleteMultipartUpload(ctx, params)
	return err
}

/*


type CompletedPart struct {

	// The base64-encoded, 32-bit CRC-32 checksum of the object. This will only be
	// present if it was uploaded with the object. When you use an API operation on an
	// object that was uploaded using multipart uploads, this value may not be a direct
	// checksum value of the full object. Instead, it's a calculation based on the
	// checksum values of each individual part. For more information about how
	// checksums are calculated with multipart uploads, see [Checking object integrity]in the Amazon S3 User
	// Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html#large-object-checksums
	ChecksumCRC32 *string

	// The base64-encoded, 32-bit CRC-32C checksum of the object. This will only be
	// present if it was uploaded with the object. When you use an API operation on an
	// object that was uploaded using multipart uploads, this value may not be a direct
	// checksum value of the full object. Instead, it's a calculation based on the
	// checksum values of each individual part. For more information about how
	// checksums are calculated with multipart uploads, see [Checking object integrity]in the Amazon S3 User
	// Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html#large-object-checksums
	ChecksumCRC32C *string

	// The base64-encoded, 160-bit SHA-1 digest of the object. This will only be
	// present if it was uploaded with the object. When you use the API operation on an
	// object that was uploaded using multipart uploads, this value may not be a direct
	// checksum value of the full object. Instead, it's a calculation based on the
	// checksum values of each individual part. For more information about how
	// checksums are calculated with multipart uploads, see [Checking object integrity]in the Amazon S3 User
	// Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html#large-object-checksums
	ChecksumSHA1 *string

	// The base64-encoded, 256-bit SHA-256 digest of the object. This will only be
	// present if it was uploaded with the object. When you use an API operation on an
	// object that was uploaded using multipart uploads, this value may not be a direct
	// checksum value of the full object. Instead, it's a calculation based on the
	// checksum values of each individual part. For more information about how
	// checksums are calculated with multipart uploads, see [Checking object integrity]in the Amazon S3 User
	// Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html#large-object-checksums
	ChecksumSHA256 *string

	// Entity tag returned when the part was uploaded.
	ETag *string

	// Part number that identifies the part. This is a positive integer between 1 and
	// 10,000.
	//
	//   - General purpose buckets - In CompleteMultipartUpload , when a additional
	//   checksum (including x-amz-checksum-crc32 , x-amz-checksum-crc32c ,
	//   x-amz-checksum-sha1 , or x-amz-checksum-sha256 ) is applied to each part, the
	//   PartNumber must start at 1 and the part numbers must be consecutive.
	//   Otherwise, Amazon S3 generates an HTTP 400 Bad Request status code and an
	//   InvalidPartOrder error code.
	//
	//   - Directory buckets - In CompleteMultipartUpload , the PartNumber must start
	//   at 1 and the part numbers must be consecutive.
	PartNumber *int32

	noSmithyDocumentSerde
}





type UploadPartInput struct {

	// The name of the bucket to which the multipart upload was initiated.
	//
	// Directory buckets - When you use this operation with a directory bucket, you
	// must use virtual-hosted-style requests in the format
	// Bucket-name.s3express-zone-id.region-code.amazonaws.com . Path-style requests
	// are not supported. Directory bucket names must be unique in the chosen Zone
	// (Availability Zone or Local Zone). Bucket names must follow the format
	// bucket-base-name--zone-id--x-s3 (for example,
	// DOC-EXAMPLE-BUCKET--usw2-az1--x-s3 ). For information about bucket naming
	// restrictions, see [Directory bucket naming rules]in the Amazon S3 User Guide.
	//
	// Access points - When you use this action with an access point, you must provide
	// the alias of the access point in place of the bucket name or specify the access
	// point ARN. When using the access point ARN, you must direct requests to the
	// access point hostname. The access point hostname takes the form
	// AccessPointName-AccountId.s3-accesspoint.Region.amazonaws.com. When using this
	// action with an access point through the Amazon Web Services SDKs, you provide
	// the access point ARN in place of the bucket name. For more information about
	// access point ARNs, see [Using access points]in the Amazon S3 User Guide.
	//
	// Access points and Object Lambda access points are not supported by directory
	// buckets.
	//
	// S3 on Outposts - When you use this action with Amazon S3 on Outposts, you must
	// direct requests to the S3 on Outposts hostname. The S3 on Outposts hostname
	// takes the form
	// AccessPointName-AccountId.outpostID.s3-outposts.Region.amazonaws.com . When you
	// use this action with S3 on Outposts through the Amazon Web Services SDKs, you
	// provide the Outposts access point ARN in place of the bucket name. For more
	// information about S3 on Outposts ARNs, see [What is S3 on Outposts?]in the Amazon S3 User Guide.
	//
	// [Directory bucket naming rules]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/directory-bucket-naming-rules.html
	// [What is S3 on Outposts?]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/S3onOutposts.html
	// [Using access points]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-access-points.html
	//
	// This member is required.
	Bucket *string

	// Object key for which the multipart upload was initiated.
	//
	// This member is required.
	Key *string

	// Part number of part being uploaded. This is a positive integer between 1 and
	// 10,000.
	//
	// This member is required.
	PartNumber *int32

	// Upload ID identifying the multipart upload whose part is being uploaded.
	//
	// This member is required.
	UploadId *string

	// Object data.
	Body io.Reader

	// Indicates the algorithm used to create the checksum for the object when you use
	// the SDK. This header will not provide any additional functionality if you don't
	// use the SDK. When you send this header, there must be a corresponding
	// x-amz-checksum or x-amz-trailer header sent. Otherwise, Amazon S3 fails the
	// request with the HTTP status code 400 Bad Request . For more information, see [Checking object integrity]
	// in the Amazon S3 User Guide.
	//
	// If you provide an individual checksum, Amazon S3 ignores any provided
	// ChecksumAlgorithm parameter.
	//
	// This checksum algorithm must be the same for all parts and it match the
	// checksum value supplied in the CreateMultipartUpload request.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumAlgorithm types.ChecksumAlgorithm

	// This header can be used as a data integrity check to verify that the data
	// received is the same data that was originally sent. This header specifies the
	// base64-encoded, 32-bit CRC-32 checksum of the object. For more information, see [Checking object integrity]
	// in the Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumCRC32 *string

	// This header can be used as a data integrity check to verify that the data
	// received is the same data that was originally sent. This header specifies the
	// base64-encoded, 32-bit CRC-32C checksum of the object. For more information, see
	// [Checking object integrity]in the Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumCRC32C *string

	// This header can be used as a data integrity check to verify that the data
	// received is the same data that was originally sent. This header specifies the
	// base64-encoded, 160-bit SHA-1 digest of the object. For more information, see [Checking object integrity]
	// in the Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumSHA1 *string

	// This header can be used as a data integrity check to verify that the data
	// received is the same data that was originally sent. This header specifies the
	// base64-encoded, 256-bit SHA-256 digest of the object. For more information, see [Checking object integrity]
	// in the Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumSHA256 *string

	// Size of the body in bytes. This parameter is useful when the size of the body
	// cannot be determined automatically.
	ContentLength *int64

	// The base64-encoded 128-bit MD5 digest of the part data. This parameter is
	// auto-populated when using the command from the CLI. This parameter is required
	// if object lock parameters are specified.
	//
	// This functionality is not supported for directory buckets.
	ContentMD5 *string

	// The account ID of the expected bucket owner. If the account ID that you provide
	// does not match the actual owner of the bucket, the request fails with the HTTP
	// status code 403 Forbidden (access denied).
	ExpectedBucketOwner *string

	// Confirms that the requester knows that they will be charged for the request.
	// Bucket owners need not specify this parameter in their requests. If either the
	// source or destination S3 bucket has Requester Pays enabled, the requester will
	// pay for corresponding charges to copy the object. For information about
	// downloading objects from Requester Pays buckets, see [Downloading Objects in Requester Pays Buckets]in the Amazon S3 User
	// Guide.
	//
	// This functionality is not supported for directory buckets.
	//
	// [Downloading Objects in Requester Pays Buckets]: https://docs.aws.amazon.com/AmazonS3/latest/dev/ObjectsinRequesterPaysBuckets.html
	RequestPayer types.RequestPayer

	// Specifies the algorithm to use when encrypting the object (for example, AES256).
	//
	// This functionality is not supported for directory buckets.
	SSECustomerAlgorithm *string

	// Specifies the customer-provided encryption key for Amazon S3 to use in
	// encrypting data. This value is used to store the object and then it is
	// discarded; Amazon S3 does not store the encryption key. The key must be
	// appropriate for use with the algorithm specified in the
	// x-amz-server-side-encryption-customer-algorithm header . This must be the same
	// encryption key specified in the initiate multipart upload request.
	//
	// This functionality is not supported for directory buckets.
	SSECustomerKey *string

	// Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321.
	// Amazon S3 uses this header for a message integrity check to ensure that the
	// encryption key was transmitted without error.
	//
	// This functionality is not supported for directory buckets.
	SSECustomerKeyMD5 *string

	noSmithyDocumentSerde
}
*/

/* Returned By CreateMultipartUpload



type CreateMultipartUploadOutput struct {

	// If the bucket has a lifecycle rule configured with an action to abort
	// incomplete multipart uploads and the prefix in the lifecycle rule matches the
	// object name in the request, the response includes this header. The header
	// indicates when the initiated multipart upload becomes eligible for an abort
	// operation. For more information, see [Aborting Incomplete Multipart Uploads Using a Bucket Lifecycle Configuration]in the Amazon S3 User Guide.
	//
	// The response also includes the x-amz-abort-rule-id header that provides the ID
	// of the lifecycle configuration rule that defines the abort action.
	//
	// This functionality is not supported for directory buckets.
	//
	// [Aborting Incomplete Multipart Uploads Using a Bucket Lifecycle Configuration]: https://docs.aws.amazon.com/AmazonS3/latest/dev/mpuoverview.html#mpu-abort-incomplete-mpu-lifecycle-config
	AbortDate *time.Time

	// This header is returned along with the x-amz-abort-date header. It identifies
	// the applicable lifecycle configuration rule that defines the action to abort
	// incomplete multipart uploads.
	//
	// This functionality is not supported for directory buckets.
	AbortRuleId *string

	// The name of the bucket to which the multipart upload was initiated. Does not
	// return the access point ARN or access point alias if used.
	//
	// Access points are not supported by directory buckets.
	Bucket *string

	// Indicates whether the multipart upload uses an S3 Bucket Key for server-side
	// encryption with Key Management Service (KMS) keys (SSE-KMS).
	BucketKeyEnabled *bool

	// The algorithm that was used to create a checksum of the object.
	ChecksumAlgorithm types.ChecksumAlgorithm

	// Object key for which the multipart upload was initiated.
	Key *string

	// If present, indicates that the requester was successfully charged for the
	// request.
	//
	// This functionality is not supported for directory buckets.
	RequestCharged types.RequestCharged

	// If server-side encryption with a customer-provided encryption key was
	// requested, the response will include this header to confirm the encryption
	// algorithm that's used.
	//
	// This functionality is not supported for directory buckets.
	SSECustomerAlgorithm *string

	// If server-side encryption with a customer-provided encryption key was
	// requested, the response will include this header to provide the round-trip
	// message integrity verification of the customer-provided encryption key.
	//
	// This functionality is not supported for directory buckets.
	SSECustomerKeyMD5 *string

	// If present, indicates the Amazon Web Services KMS Encryption Context to use for
	// object encryption. The value of this header is a Base64-encoded string of a
	// UTF-8 encoded JSON, which contains the encryption context as key-value pairs.
	SSEKMSEncryptionContext *string

	// If present, indicates the ID of the KMS key that was used for object encryption.
	SSEKMSKeyId *string

	// The server-side encryption algorithm used when you store this object in Amazon
	// S3 (for example, AES256 , aws:kms ).
	ServerSideEncryption types.ServerSideEncryption

	// ID for the initiated multipart upload.
	UploadId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}
*/
