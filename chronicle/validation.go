package chronicle

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	googleoauth "golang.org/x/oauth2/google"
)

func validateCredentials(v interface{}, k cty.Path) diag.Diagnostics {
	creds := v.(string)

	if _, err := os.Stat(creds); err == nil {
		return nil
	}
	if _, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(creds)); err != nil {
		return diag.FromErr(fmt.Errorf("JSON credentials are not valid: %s", err))
	}

	return nil
}

func validateRegion(v interface{}, k cty.Path) diag.Diagnostics {
	region := v.(string)

	if !isValidRegion(region) {
		return diag.FromErr(fmt.Errorf("region %s not valid, valid regions are: %s", region, chronicle.Regions))
	}

	return nil
}

func isValidRegion(region string) bool {
	return contains(chronicle.Regions, region)
}

func validateFeedS3SourceDeleteOption(v interface{}, k cty.Path) diag.Diagnostics {
	deletionOptions := []string{FeedS3SourceDeleteOptionDeletionNever, FeedS3SourceDeleteOptionDeletionOnSuccess, FeedS3SourceDeleteOptionDeletionOnSuccessFilesOnly}
	option := v.(string)
	if !contains(deletionOptions, option) {
		return diag.FromErr(fmt.Errorf("source deletion option %s not valid, valid options are: %s", option, deletionOptions))
	}

	return nil
}

func validateFeedGCSSourceDeleteOption(v interface{}, k cty.Path) diag.Diagnostics {
	deletionOptions := []string{FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionNever, FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionOnSuccess,
		FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionOnSuccessFilesOnly}
	option := v.(string)
	if !contains(deletionOptions, option) {
		return diag.FromErr(fmt.Errorf("source deletion option %s not valid, valid options are: %s", option, deletionOptions))
	}

	return nil
}

func validateFeedAzureBlobStoreSourceDeleteOption(v interface{}, k cty.Path) diag.Diagnostics {
	deletionOptions := []string{FeedAzureBlobStoreSourceDeleteOptionDeletionNever}
	option := v.(string)
	if !contains(deletionOptions, option) {
		return diag.FromErr(fmt.Errorf("source deletion option %s not valid, valid options are: %s", option, deletionOptions))
	}

	return nil
}

func validateFeedS3SourceType(v interface{}, k cty.Path) diag.Diagnostics {
	sourceTypes := []string{FeedS3SourceTypeFiles, FeedS3SourceTypeFolders, FeedS3SourceTypeFoldersRecursive}
	sourceType := v.(string)
	if !contains(sourceTypes, sourceType) {
		return diag.FromErr(fmt.Errorf("source type %s not valid, valid types are: %s", sourceType, sourceTypes))
	}

	return nil
}

func validateFeedGCSSourceType(v interface{}, k cty.Path) diag.Diagnostics {
	sourceTypes := []string{FeedGoogleCloudStorageBucketSourceTypeFiles, FeedGoogleCloudStorageBucketSourceTypeFolders, FeedGoogleCloudStorageBucketSourceTypeFoldersRecursive}
	sourceType := v.(string)
	if !contains(sourceTypes, sourceType) {
		return diag.FromErr(fmt.Errorf("source type %s not valid, valid types are: %s", sourceType, sourceTypes))
	}

	return nil
}

func validateFeedAzureBlobStoreSourceType(v interface{}, k cty.Path) diag.Diagnostics {
	sourceTypes := []string{FeedAzureBlobStoreSourceTypeFiles, FeedAzureBlobStoreSourceTypeFolders, FeedAzureBlobStoreSourceTypeFoldersRecursive}
	sourceType := v.(string)
	if !contains(sourceTypes, sourceType) {
		return diag.FromErr(fmt.Errorf("source type %s not valid, valid types are: %s", sourceType, sourceTypes))
	}

	return nil
}

func validateRuleText(v interface{}, k cty.Path) diag.Diagnostics {
	ruleText := v.(string)
	if !strings.HasSuffix(ruleText, "\n") {
		return diag.FromErr(fmt.Errorf("rule_text %s not valid, it must end with new line", ruleText))
	}

	return nil
}

func validateCustomEndpoint(v interface{}, k cty.Path) diag.Diagnostics {
	u := v.(string)
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%q cannot be validated", u))
	}

	return nil
}

func validateSubjectType(v interface{}, k cty.Path) diag.Diagnostics {
	subjectTypes := []string{RBACSubjectTypeAnalyst, RBACSubjectTypeIDPGroup}
	subjectType := v.(string)
	if !contains(subjectTypes, subjectType) {
		return diag.FromErr(fmt.Errorf("subject type %s not valid, valid types are: %s", subjectType, subjectTypes))
	}

	return nil
}

func validateRegexp(re string) schema.SchemaValidateDiagFunc {
	return func(v interface{}, k cty.Path) diag.Diagnostics {
		value := v.(string)
		if !regexp.MustCompile(re).MatchString(value) {
			return diag.FromErr(fmt.Errorf("%q (%q) doesn't match regexp %q", k, value, re))
		}

		return nil
	}
}

func validateGCSURI(v interface{}, k cty.Path) diag.Diagnostics {
	reg := `^gs:\/\/([a-z0-9._\-]+\/?)+$`
	return validateRegexp(reg)(v, k)
}

func validateAWSAccessKeyID(v interface{}, k cty.Path) diag.Diagnostics {
	reg := `^[A-Z0-9]{20}$`
	return validateRegexp(reg)(v, k)
}
func validateAWSSecretAccessKey(v interface{}, k cty.Path) diag.Diagnostics {
	reg := `^[A-Za-z0-9\/+=]{40}$`
	return validateRegexp(reg)(v, k)
}

func validateAWSAccountID(v interface{}, k cty.Path) diag.Diagnostics {
	reg := `^\d{12}$`
	return validateRegexp(reg)(v, k)
}

func validateThinkstCanaryHostname(v interface{}, k cty.Path) diag.Diagnostics {
	reg := `^.*\.canary\.tools$`
	return validateRegexp(reg)(v, k)
}

func validateUUID(v interface{}, k cty.Path) diag.Diagnostics {
	u := v.(string)
	_, err := uuid.Parse(u)
	if err != nil {
		return diag.FromErr(fmt.Errorf("uuid %s not valid", u))
	}

	return nil
}

func validateFeedMicrosoftOffice365ManagementActivityContentType(v interface{}, k cty.Path) diag.Diagnostics {
	contentTypes := []string{FeedMicrosoftOffice365ManagementActivityContentTypeAuditAzureActiveDirectory,
		FeedMicrosoftOffice365ManagementActivityContentTypeAuditExchange,
		FeedMicrosoftOffice365ManagementActivityContentTypeAuditSharePoint,
		FeedMicrosoftOffice365ManagementActivityContentTypeAuditGeneral,
		FeedMicrosoftOffice365ManagementActivityContentTypeDPLAll}
	contentType := v.(string)
	if !contains(contentTypes, contentType) {
		return diag.FromErr(fmt.Errorf("conten type %s not valid, valid types are: %s", contentType, contentTypes))
	}
	return nil
}

func validateReferenceListContentType(v interface{}, k cty.Path) diag.Diagnostics {
	contentTypes := []string{string(chronicle.ReferenceListContentTypeCIDR),
		string(chronicle.ReferenceListContentTypeREGEX),
		string(chronicle.ReferenceListContentTypeDefault)}
	contentType := v.(string)
	if !contains(contentTypes, contentType) {
		return diag.FromErr(fmt.Errorf("conten type %s not valid, valid types are: %s", contentType, contentTypes))
	}
	return nil
}
