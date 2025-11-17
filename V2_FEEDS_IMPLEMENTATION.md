# V2 Feed Implementation Summary

This document summarizes the implementation of V2 feed support for the Chronicle Terraform Provider.

## Overview

Four new V2 feed resources have been implemented as **standalone, independent resources** separate from their V1 counterparts:

1. **chronicle_feed_amazon_sqs_v2** - Amazon SQS V2 feed
2. **chronicle_feed_amazon_s3_v2** - Amazon S3 V2 feed
3. **chronicle_feed_google_cloud_storage_v2** - Google Cloud Storage V2 feed
4. **chronicle_feed_azure_blobstore_v2** - Azure Blob Storage V2 feed

All V2 feeds use the **Google Cloud Storage Transfer Service** for improved ingestion performance.

## Key Differences from V1

### Common V2 Features

All V2 feeds share these characteristics:

- **Transfer Service**: Use Google Cloud Storage Transfer Service instead of direct polling
- **Source Deletion Options**: Only 2 options (vs V1's 3):
  - `NEVER`: Never delete files from source
  - `ON_SUCCESS`: Delete files and empty directories after successful ingestion
- **Max Lookback Days**: Configurable file age limit (default: 180 days)
- **Improved Performance**: Better ingestion latency and throughput

### Feed-Specific Details

#### AMAZON_SQS_V2
- **Feed Source Type**: `AMAZON_SQS_V2`
- **Configuration Property**: `amazonSqsV2Settings`
- **New Fields**:
  - `s3_uri` - S3 bucket URI
  - `queue_name` - SQS queue name
  - `max_lookback_days` - File age limit
- **Authentication**: Dual authentication (SQS + optional S3)
- **Resource Name**: `chronicle_feed_amazon_sqs_v2`

#### AMAZON_S3_V2
- **Feed Source Type**: `AMAZON_S3_V2`
- **Configuration Property**: `amazonS3V2Settings`
- **Removed Fields**: `sourceType` (no longer needed)
- **New Fields**:
  - `max_lookback_days` - File age limit
- **Authentication**: AWS access keys with region
- **Resource Name**: `chronicle_feed_amazon_s3_v2`

#### GOOGLE_CLOUD_STORAGE_V2
- **Feed Source Type**: `GOOGLE_CLOUD_STORAGE_V2`
- **Configuration Property**: `gcsV2Settings`
- **Removed Fields**: `sourceType` (no longer needed)
- **New Fields**:
  - `max_lookback_days` - File age limit
- **Authentication**: Service account (no explicit credentials in config)
- **Resource Name**: `chronicle_feed_google_cloud_storage_v2`

#### AZURE_BLOBSTORE_V2
- **Feed Source Type**: `AZURE_BLOBSTORE_V2`
- **Configuration Property**: `azureBlobStoreV2Settings`
- **Removed Fields**: `sourceType` (no longer needed)
- **New Fields**:
  - `max_lookback_days` - File age limit
- **Authentication**: Shared Key OR SAS Token (mutually exclusive)
- **Resource Name**: `chronicle_feed_azure_blobstore_v2`

## Implementation Files

### Client Layer (API Interaction)

| File | Purpose |
|------|---------|
| `client/feed.go` | Updated with V2 feed source type constants |
| `client/feed_amazon_sqs_v2.go` | SQS V2 configuration struct and methods |
| `client/feed_amazon_s3_v2.go` | S3 V2 configuration struct and methods |
| `client/feed_google_cloud_storage_v2.go` | GCS V2 configuration struct and methods |
| `client/feed_azure_blobstore_v2.go` | Azure V2 configuration struct and methods |

### Resource Layer (Terraform Integration)

| File | Purpose |
|------|---------|
| `chronicle/resource_feed_amazon_sqs_v2.go` | SQS V2 Terraform resource |
| `chronicle/resource_feed_amazon_s3_v2.go` | S3 V2 Terraform resource |
| `chronicle/resource_feed_google_cloud_storage_v2.go` | GCS V2 Terraform resource |
| `chronicle/resource_feed_azure_blobstore_v2.go` | Azure V2 Terraform resource |
| `chronicle/validation.go` | Added `validateFeedV2SourceDeleteOption` |
| `chronicle/provider.go` | Registered all 4 V2 resources |

### Examples

| Directory | Contents |
|-----------|----------|
| `examples/resources/chronicle_feed_amazon_sqs_v2/` | SQS V2 example configuration |
| `examples/resources/chronicle_feed_amazon_s3_v2/` | S3 V2 example configuration |
| `examples/resources/chronicle_feed_google_cloud_storage_v2/` | GCS V2 example configuration |
| `examples/resources/chronicle_feed_azure_blobstore_v2/` | Azure V2 example configuration |

## Resource Schema

### Common Fields (All V2 Feeds)

```hcl
resource "chronicle_feed_<type>_v2" "example" {
  display_name = "Feed Name"          # Required
  log_type     = "LOG_TYPE"           # Required
  enabled      = true                  # Required
  namespace    = "namespace"           # Optional
  labels       = { key = "value" }     # Optional

  details {
    # Feed-specific configuration
  }
}
```

### Validation Rules

- **source_delete_options**: Must be `NEVER` or `ON_SUCCESS`
- **max_lookback_days**: Integer, default 180
- **AWS Access Key ID**: 20-character alphanumeric string
- **AWS Secret Access Key**: 40-character string (alphanumeric + special chars)
- **AWS Account ID**: 12-digit number
- **GCS URI**: Must match pattern `gs://bucket/path/`
- **Azure Authentication**: Must specify either `shared_key` OR `sas_token`, not both

## Testing Considerations

### Unit Tests Needed

For each V2 feed type, implement:

1. **Expand/Flatten Tests**:
   - Test conversion from Terraform config to API request
   - Test conversion from API response to Terraform state
   - Test handling of optional fields
   - Test nil/empty value handling

2. **Validation Tests**:
   - Test source delete option validation
   - Test AWS credential format validation
   - Test GCS URI format validation
   - Test Azure auth mutual exclusivity

### Acceptance Tests Needed

For each V2 feed type, implement:

1. **Basic CRUD**:
   - Create feed with minimum required fields
   - Read feed and verify all fields
   - Update feed configuration
   - Delete feed
   - Import existing feed by ID

2. **Advanced Scenarios**:
   - Create with all optional fields
   - Update source_delete_options
   - Update max_lookback_days
   - Enable/disable feed
   - Update labels and namespace

3. **Error Scenarios**:
   - Invalid credentials
   - Invalid source_delete_options
   - Invalid max_lookback_days
   - API errors (404, 403, etc.)

### Test Environment Requirements

- Valid AWS credentials for S3/SQS testing
- Valid GCP service account for GCS testing
- Valid Azure credentials for Blob Storage testing
- Chronicle/Google Security Operations API access
- Test buckets/queues/containers in each cloud provider

## Migration from V1

**Important**: There is **NO automated migration path** from V1 to V2 feeds. Users must:

1. Create new V2 feed resources separately
2. Test V2 feeds in parallel with V1
3. Manually cutover when ready
4. Destroy V1 feeds when no longer needed

V1 resources will eventually be deprecated, but that is a separate effort.

## Usage Examples

See the `/examples/resources/chronicle_feed_*_v2/` directories for complete examples of each feed type.

### Quick Start: Amazon S3 V2

```hcl
resource "chronicle_feed_amazon_s3_v2" "cloudtrail" {
  display_name = "CloudTrail Logs"
  log_type     = "AWS_CLOUDTRAIL"
  enabled      = true

  details {
    s3_uri                = "s3://my-logs/cloudtrail/"
    source_delete_options = "ON_SUCCESS"
    max_lookback_days     = 180

    authentication {
      region            = "us-east-1"
      access_key_id     = var.aws_access_key_id
      secret_access_key = var.aws_secret_access_key
    }
  }
}
```

## API Reference

The V2 feeds use the Chronicle Feed Management API with the following endpoints:

- **Create**: `POST /v1/feeds`
- **Read**: `GET /v1/feeds/{id}`
- **Update**: `PATCH /v1/feeds/{id}`
- **Delete**: `DELETE /v1/feeds/{id}`
- **Enable/Disable**: `POST /v1/feeds/{id}:enable` or `POST /v1/feeds/{id}:disable`

Request bodies include the feed configuration under the appropriate property key:
- SQS V2: `amazonSqsV2Settings`
- S3 V2: `amazonS3V2Settings`
- GCS V2: `gcsV2Settings`
- Azure V2: `azureBlobStoreV2Settings`

## Known Limitations

1. **No V1 Migration**: Users must manually recreate feeds
2. **Service Account Setup**: GCS V2 requires manual service account configuration
3. **Credential Visibility**: Credentials are not returned in API read operations
4. **Source Delete Options**: Only 2 options vs V1's 3

## Build and Installation

```bash
# Build the provider
go build -o terraform-provider-chronicle

# Install locally
make install

# Run tests
make test

# Run acceptance tests (requires API credentials)
make testacc
```

## Next Steps

1. **Write Tests**: Implement unit and acceptance tests
2. **Generate Docs**: Run `tfplugindocs` to auto-generate documentation
3. **Update Changelog**: Document new V2 feed resources
4. **Integration Testing**: Test with real Chronicle API
5. **User Documentation**: Create migration guides and best practices

## References

- [Chronicle Feed Management API](https://cloud.google.com/chronicle/docs/reference/feed-management-api)
- [Existing V1 Implementation](./chronicle/resource_feed_amazon_s3.go)
- [Provider Schema](./chronicle/provider.go)
