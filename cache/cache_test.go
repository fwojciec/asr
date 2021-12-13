package cache_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/fwojciec/asr"
	"github.com/fwojciec/asr/cache"
)

func TestAddService(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		data []*asr.Service
		exp  *cache.Cache
	}{
		{
			name: "basic service",
			data: []*asr.Service{{
				Name:         "test_name",
				Prefix:       "test_prefix",
				ConfigDocURL: "test_config_doc_url",
				APIDocURL:    "test_api_doc_url",
				IAMDocURL:    "test_iam_doc_url",
			}},
			exp: &cache.Cache{
				ServiceByID: map[string]*cache.Service{"test_prefix": {
					ID:           "test_prefix",
					Name:         "test_name",
					Prefix:       "test_prefix",
					ConfigDocURL: "test_config_doc_url",
					APIDocURL:    "test_api_doc_url",
					IAMDocURL:    "test_iam_doc_url",
				}},
				SortedServiceIDs: []string{"test_prefix"},
			},
		},
		{
			name: "service action",
			data: []*asr.Service{{
				Prefix: "test_prefix",
				Actions: []asr.Action{{
					Name:        "test_action_name",
					DocURL:      "test_doc_url",
					Description: "test_description",
					AccessLevel: "test_access_level",
					ResourceTypes: []asr.ActionResourceType{{
						Name:     "test_resource_type",
						Required: true,
					}},
					ConditionKeys:    []string{"test_condition_key"},
					DependentActions: []string{"test_dependent_action"},
				}},
			}},
			exp: &cache.Cache{
				ServiceByID: map[string]*cache.Service{"test_prefix": {
					ID:      "test_prefix",
					Prefix:  "test_prefix",
					Actions: []string{"test_prefix:test_action_name"},
				}},
				ActionByID: map[string]*cache.Action{"test_prefix:test_action_name": {
					ID:            "test_prefix:test_action_name",
					Name:          "test_action_name",
					DocURL:        "test_doc_url",
					Description:   "test_description",
					AccessLevel:   "test_access_level",
					ConditionKeys: []string{"test_condition_key"},
					ResourceTypes: []cache.ActionResourceType{{
						ResourceType: "test_prefix:test_resource_type",
						Required:     true,
					}},
					DependentActions: []string{"test_dependent_action"},
					Service:          "test_prefix",
				}},
				ActionIDsByConditionKeyID:    map[string][]string{"test_condition_key": {"test_prefix:test_action_name"}},
				ActionIDsByResourceTypeID:    map[string][]string{"test_prefix:test_resource_type": {"test_prefix:test_action_name"}},
				ActionIDsByDependentActionID: map[string][]string{"test_dependent_action": {"test_prefix:test_action_name"}},
				SortedServiceIDs:             []string{"test_prefix"},
				SortedActionIDs:              []string{"test_prefix:test_action_name"},
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := cache.NewCache(tc.data)
			equals(t, tc.exp, res)
		})
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
