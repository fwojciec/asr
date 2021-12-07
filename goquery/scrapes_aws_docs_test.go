package goquery_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/fwojciec/asr"
	"github.com/fwojciec/asr/goquery"
)

var testFileCache = map[string][]byte{}
var testFileCacheLock sync.Mutex

func openTestFile(path string) (io.Reader, error) {
	testFileCacheLock.Lock()
	defer testFileCacheLock.Unlock()
	if c := testFileCache[path]; c != nil {
		return bytes.NewBuffer(c), nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	testFileCache[path] = b
	return bytes.NewBuffer(b), nil
}

type mockGetter struct{}

func (m *mockGetter) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	f, err := openTestFile(url)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(f), nil
}

func TestScrapesPreamblePrefix(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		path string
		exp  string
	}{
		{"testdata/alexaforbusiness.html", "a4b"},
		{"testdata/amazonapigateway.html", "execute-api"},
		{"testdata/amazonec2autoscaling.html", "autoscaling"},
		{"testdata/awsaccountmanagement.html", "account"},
		{"testdata/awsamplify.html", "amplify"},
	} {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			subject := goquery.NewScrapesAWSDocs(&mockGetter{})
			res, err := subject.Scrape(context.Background(), []*asr.TOCEntry{{URL: tc.path}})
			ok(t, err)
			equals(t, tc.exp, res[0].Prefix)
		})
	}
}

func TestScrapesPreambleConfigDocURL(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		path string
		exp  string
	}{
		{"testdata/awsaccountmanagement.html", "https://docs.aws.amazon.com/accounts/latest/reference/"},
		{"testdata/alexaforbusiness.html", ""},
		{"testdata/amazonapigateway.html", "https://docs.aws.amazon.com/apigateway/latest/developerguide/"},
		{"testdata/amazonec2autoscaling.html", "https://docs.aws.amazon.com/autoscaling/latest/userguide/"},
	} {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			subject := goquery.NewScrapesAWSDocs(&mockGetter{})
			res, err := subject.Scrape(context.Background(), []*asr.TOCEntry{{URL: tc.path}})
			ok(t, err)
			equals(t, tc.exp, res[0].ConfigDocURL)
		})
	}
}

func TestScrapesPreambleAPIDocURL(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		path string
		exp  string
	}{
		{"testdata/awsaccountmanagement.html", "https://docs.aws.amazon.com/accounts/latest/reference/api"},
		{"testdata/alexaforbusiness.html", "https://docs.aws.amazon.com/a4b/latest/APIReference/"},
		{"testdata/amazonapigateway.html", "https://docs.aws.amazon.com/apigateway/api-reference/"},
		{"testdata/amazonec2autoscaling.html", "https://docs.aws.amazon.com/AutoScaling/latest/APIReference/"},
	} {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			subject := goquery.NewScrapesAWSDocs(&mockGetter{})
			res, err := subject.Scrape(context.Background(), []*asr.TOCEntry{{URL: tc.path}})
			ok(t, err)
			equals(t, tc.exp, res[0].APIDocURL)
		})
	}
}

func TestScrapesPreambleIAMDocURL(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		path string
		exp  string
	}{
		{"testdata/awsaccountmanagement.html", "${UserGuideDocPage}security_iam_service-with-iam.html"},
		{"testdata/alexaforbusiness.html", ""},
		{"testdata/amazonapigateway.html", "https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-control-access-to-api.html"},
		{"testdata/amazonec2autoscaling.html", "https://docs.aws.amazon.com/autoscaling/latest/userguide/IAM.html"},
	} {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			subject := goquery.NewScrapesAWSDocs(&mockGetter{})
			res, err := subject.Scrape(context.Background(), []*asr.TOCEntry{{URL: tc.path}})
			ok(t, err)
			equals(t, tc.exp, res[0].IAMDocURL)
		})
	}
}

func TestScrapesResourceTypes(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		path string
		exp  []asr.ResourceType
	}{
		{
			path: "testdata/awsaccountmanagement.html",
			exp: []asr.ResourceType{
				{
					Name:       "account",
					DocURL:     "${UserGuideDocPage}security_iam_service-with-iam.html#security_iam_service-with-iam-id-based-policies-resources",
					ARNPattern: "arn:${Partition}:account::${Account}:account",
				},
				{
					Name:       "accountInOrganization",
					DocURL:     "${UserGuideDocPage}security_iam_service-with-iam.html#security_iam_service-with-iam-id-based-policies-resources",
					ARNPattern: "arn:${Partition}:account::${ManagementAccountId}:account/o-${OrganizationId}/${MemberAccountId}",
				},
			},
		},
		{
			path: "testdata/amazonapigateway.html",
			exp: []asr.ResourceType{
				{
					Name:       "execute-api-general",
					DocURL:     "",
					ARNPattern: "arn:${Partition}:execute-api:${Region}:${Account}:${ApiId}/${Stage}/${Method}/${ApiSpecificResourcePath}",
				},
			},
		},
		{
			path: "testdata/amazonec2autoscaling.html",
			exp: []asr.ResourceType{
				{
					Name:          "autoScalingGroup",
					DocURL:        "https://docs.aws.amazon.com/autoscaling/latest/userguide/control-access-using-iam.html#policy-auto-scaling-resources",
					ARNPattern:    "arn:${Partition}:autoscaling:${Region}:${Account}:autoScalingGroup:${GroupId}:autoScalingGroupName/${GroupFriendlyName}",
					ConditionKeys: []string{"autoscaling:ResourceTag/${TagKey}", "aws:ResourceTag/${TagKey}"},
				},
				{
					Name:       "launchConfiguration",
					DocURL:     "https://docs.aws.amazon.com/autoscaling/latest/userguide/control-access-using-iam.html#policy-auto-scaling-resources",
					ARNPattern: "arn:${Partition}:autoscaling:${Region}:${Account}:launchConfiguration:${Id}:launchConfigurationName/${LaunchConfigurationName}",
				},
			},
		},
	} {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			subject := goquery.NewScrapesAWSDocs(&mockGetter{})
			res, err := subject.Scrape(context.Background(), []*asr.TOCEntry{{URL: tc.path}})
			ok(t, err)
			equals(t, tc.exp, res[0].ResourceTypes)
		})
	}
}

func TestScrapesConditionKeys(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		path string
		exp  []asr.ConditionKey
	}{
		{
			path: "testdata/awsaccountmanagement.html",
			exp: []asr.ConditionKey{
				{Name: "account:AccountResourceOrgPaths", DocURL: "${UserGuideDocPage}security_iam_service-with-iam.html#security_iam_service-with-iam-id-based-policies-conditionkeys", Description: "Filters access by the resource path for an account in an organization", Type: "ArrayOfString"},
				{Name: "account:AccountResourceOrgTags/${TagKey}", DocURL: "${UserGuideDocPage}security_iam_service-with-iam.html#security_iam_service-with-iam-id-based-policies-conditionkeys", Description: "Filters access by resource tags for an account in an organization", Type: "ArrayOfString"},
				{Name: "account:AlternateContactTypes", DocURL: "${UserGuideDocPage}security_iam_service-with-iam.html#security_iam_service-with-iam-id-based-policies-conditionkeys", Description: "Filters access by alternate contact types", Type: "ArrayOfString"},
				{Name: "account:TargetRegion", DocURL: "${UserGuideDocPage}security_iam_service-with-iam.html#security_iam_service-with-iam-id-based-policies-conditionkeys", Description: "Filters access by a list of Regions. Enables or disables all the Regions specified here", Type: "String"},
			},
		},
		{
			path: "testdata/amazonapigateway.html",
		},
		{
			path: "testdata/awsamplify.html",
			exp: []asr.ConditionKey{
				{Name: "aws:RequestTag/${TagKey}", Type: "String"},
				{Name: "aws:ResourceTag/${TagKey}", Type: "String"},
				{Name: "aws:TagKeys", Type: "String"},
			},
		},
	} {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			subject := goquery.NewScrapesAWSDocs(&mockGetter{})
			res, err := subject.Scrape(context.Background(), []*asr.TOCEntry{{URL: tc.path}})
			ok(t, err)
			equals(t, tc.exp, res[0].ConditionKeys)
		})
	}
}

func TestScrapesActions(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		path string
		exp  []asr.Action
	}{
		{
			path: "testdata/awsaccountmanagement.html",
			exp: []asr.Action{
				{
					Name:          "DeleteAlternateContact",
					DocURL:        "${APIReferenceDocPage}API_DeleteAlternateContact.html",
					Description:   "Grants permission to delete the alternate contacts for an account",
					AccessLevel:   "Write",
					ResourceTypes: []asr.ActionResourceType{{Name: "account", Required: false}, {Name: "accountInOrganization", Required: false}},
				},
				{
					Name:          "DisableRegion",
					DocURL:        "https://docs.aws.amazon.com/general/latest/gr/rande-manage.html",
					Description:   "Grants permission to disable use of a Region",
					AccessLevel:   "Write",
					ConditionKeys: []string{"account:TargetRegion"},
				},
				{
					Name:          "EnableRegion",
					DocURL:        "https://docs.aws.amazon.com/general/latest/gr/rande-manage.html",
					Description:   "Grants permission to enable use of a Region",
					AccessLevel:   "Write",
					ConditionKeys: []string{"account:TargetRegion"},
				},
				{
					Name:          "GetAlternateContact",
					DocURL:        "${APIReferenceDocPage}API_GetAlternateContact.html",
					Description:   "Grants permission to retrieve the alternate contacts for an account",
					AccessLevel:   "Read",
					ResourceTypes: []asr.ActionResourceType{{Name: "account", Required: false}, {Name: "accountInOrganization", Required: false}},
				},
				{
					Name:        "ListRegions",
					DocURL:      "https://docs.aws.amazon.com/general/latest/gr/rande-manage.html",
					Description: "Grants permission to list the available Regions",
					AccessLevel: "List",
				},
				{
					Name:          "PutAlternateContact",
					DocURL:        "${APIReferenceDocPage}API_PutAlternateContact.html",
					Description:   "Grants permission to modify the alternate contacts for an account",
					AccessLevel:   "Write",
					ResourceTypes: []asr.ActionResourceType{{Name: "account", Required: false}, {Name: "accountInOrganization", Required: false}},
				},
			},
		},
		{
			path: "testdata/amazonapigateway.html",
			exp: []asr.Action{
				{
					Name:          "InvalidateCache",
					DocURL:        "https://docs.aws.amazon.com/apigateway/api-reference/api-gateway-caching.html",
					Description:   "Used to invalidate API cache upon a client request",
					AccessLevel:   "Write",
					ResourceTypes: []asr.ActionResourceType{{Name: "execute-api-general", Required: true}},
				},
				{
					Name:          "Invoke",
					DocURL:        "https://docs.aws.amazon.com/apigateway/api-reference/how-to-call-api.html",
					Description:   "Used to invoke an API upon a client request",
					AccessLevel:   "Write",
					ResourceTypes: []asr.ActionResourceType{{Name: "execute-api-general", Required: true}},
				},
				{
					Name:          "ManageConnections",
					DocURL:        "https://docs.aws.amazon.com/apigateway/api-reference/apigateway-websocket-control-access-iam.html",
					Description:   "ManageConnections controls access to the @connections API",
					AccessLevel:   "Write",
					ResourceTypes: []asr.ActionResourceType{{Name: "execute-api-general", Required: true}},
				},
			},
		},
		{
			path: "testdata/awsbackupstorage.html",
			exp: []asr.Action{
				{
					Name:        "MountCapsule",
					DocURL:      "https://docs.aws.amazon.com/aws-backup/latest/devguide/API_CreateBackupVault.html",
					Description: "Associates a KMS key to a backup vault",
					AccessLevel: "Write",
				},
			},
		},
		{
			path: "testdata/awscloudwatchrum.html",
			exp: []asr.Action{
				{
					Name:             "CreateAppMonitor",
					DocURL:           "https://docs.aws.amazon.com/rum/latest/APIReference/API_CreateAppMonitor.html",
					Description:      "Grants permission to create appMonitor metadata",
					AccessLevel:      "Write",
					ResourceTypes:    []asr.ActionResourceType{{Name: "AppMonitorResource", Required: true}},
					ConditionKeys:    []string{"aws:RequestTag/${TagKey}", "aws:TagKeys"},
					DependentActions: []string{"iam:CreateServiceLinkedRole", "iam:GetRole"},
				},
				{
					Name:          "DeleteAppMonitor",
					DocURL:        "https://docs.aws.amazon.com/rum/latest/APIReference/API_DeleteAppMonitor.html",
					Description:   "Grants permission to delete appMonitor metadata",
					AccessLevel:   "Write",
					ResourceTypes: []asr.ActionResourceType{{Name: "AppMonitorResource", Required: true}},
				},
				{
					Name:          "GetAppMonitor",
					DocURL:        "https://docs.aws.amazon.com/rum/latest/APIReference/API_GetAppMonitor.html",
					Description:   "Grants permission to get appMonitor metadata",
					AccessLevel:   "Read",
					ResourceTypes: []asr.ActionResourceType{{Name: "AppMonitorResource", Required: true}},
				},
				{
					Name:          "GetAppMonitorData",
					DocURL:        "https://docs.aws.amazon.com/rum/latest/APIReference/API_GetAppMonitorData.html",
					Description:   "Grants permission to get appMonitor data",
					AccessLevel:   "Read",
					ResourceTypes: []asr.ActionResourceType{{Name: "AppMonitorResource", Required: true}},
				},
				{
					Name:        "ListAppMonitors",
					DocURL:      "https://docs.aws.amazon.com/rum/latest/APIReference/API_ListAppMonitors.html",
					Description: "Grants permission to list appMonitors metadata",
					AccessLevel: "List",
				},
				{
					Name:        "ListTagsForResource",
					DocURL:      "https://docs.aws.amazon.com/rum/latest/APIReference/API_ListTagsForResource.html",
					Description: "Grants permission to list tags for resources",
					AccessLevel: "Read",
				},
				{
					Name:        "PutRumEvents",
					DocURL:      "https://docs.aws.amazon.com/rum/latest/APIReference/API_PutRumEvents.html",
					Description: "Grants permission to put RUM events for appmonitor",
					AccessLevel: "Write",
				},
				{
					Name:        "TagResource",
					DocURL:      "https://docs.aws.amazon.com/rum/latest/APIReference/API_TagResource.html",
					Description: "Grants permission to tag resources",
					AccessLevel: "Tagging",
				},
				{
					Name:        "UntagResource",
					DocURL:      "https://docs.aws.amazon.com/rum/latest/APIReference/API_UntagResource.html",
					Description: "Grants permission to untag resources",
					AccessLevel: "Tagging",
				},
				{
					Name:             "UpdateAppMonitor",
					DocURL:           "https://docs.aws.amazon.com/rum/latest/APIReference/API_UpdateAppMonitor.html",
					Description:      "Grants permission to update appmonitor metadata",
					AccessLevel:      "Write",
					ResourceTypes:    []asr.ActionResourceType{{Name: "AppMonitorResource", Required: true}},
					DependentActions: []string{"iam:CreateServiceLinkedRole", "iam:GetRole"},
				},
			},
		},
	} {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			subject := goquery.NewScrapesAWSDocs(&mockGetter{})
			res, err := subject.Scrape(context.Background(), []*asr.TOCEntry{{URL: tc.path}})
			ok(t, err)
			equals(t, tc.exp, res[0].Actions)
		})
	}
}
