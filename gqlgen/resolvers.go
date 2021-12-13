package gqlgen

import (
	"context"
	"strings"

	"github.com/fwojciec/asr/cache"
)

type Resolver struct {
	Cache *cache.Cache
}

func (r *actionResolver) ConditionKeys(ctx context.Context, obj *cache.Action) ([]*cache.ConditionKey, error) {
	res := make([]*cache.ConditionKey, len(obj.ConditionKeys))
	for i, ck := range obj.ConditionKeys {
		res[i] = r.Cache.ConditionKeyByID[ck]
	}
	return res, nil
}

func (r *actionResolver) DependentActions(ctx context.Context, obj *cache.Action) ([]*cache.Action, error) {
	res := make([]*cache.Action, len(obj.DependentActions))
	for i, ck := range obj.DependentActions {
		res[i] = r.Cache.ActionByID[ck]
	}
	return res, nil
}

func (r *actionResolver) Service(ctx context.Context, obj *cache.Action) (*cache.Service, error) {
	return r.Cache.ServiceByID[obj.Service], nil
}

func (r *actionResourceTypeResolver) ResourceType(ctx context.Context, obj *cache.ActionResourceType) (*cache.ResourceType, error) {
	return r.Cache.ResourceTypeByID[obj.ResourceType], nil
}

func (r *conditionKeyResolver) Actions(ctx context.Context, obj *cache.ConditionKey) ([]*cache.Action, error) {
	aIDs := r.Cache.ActionIDsByConditionKeyID[obj.ID]
	res := make([]*cache.Action, len(aIDs))
	for i, aID := range aIDs {
		res[i] = r.Cache.ActionByID[aID]
	}
	return res, nil
}

func (r *conditionKeyResolver) ResourceTypes(ctx context.Context, obj *cache.ConditionKey) ([]*cache.ResourceType, error) {
	rtIDs := r.Cache.ResourceTypeIDsByConditionKeyID[obj.ID]
	res := make([]*cache.ResourceType, len(rtIDs))
	for i, rtID := range rtIDs {
		res[i] = r.Cache.ResourceTypeByID[rtID]
	}
	return res, nil
}

func (r *conditionKeyResolver) Services(ctx context.Context, obj *cache.ConditionKey) ([]*cache.Service, error) {
	sIDs := r.Cache.ServiceIDsByConditionKeyID[obj.ID]
	res := make([]*cache.Service, len(sIDs))
	for i, sID := range sIDs {
		res[i] = r.Cache.ServiceByID[sID]
	}
	return res, nil
}

func (r *queryResolver) Services(ctx context.Context, filter *string) ([]*cache.Service, error) {
	if filter != nil {
		return r.servicesWithFilter(*filter), nil
	}
	return r.allServices(), nil
}

func (r *queryResolver) servicesWithFilter(filter string) []*cache.Service {
	var ids []string
	for _, id := range r.Cache.SortedServiceIDs {
		if !strings.Contains(id, filter) {
			continue
		}
		ids = append(ids, id)
	}
	res := make([]*cache.Service, len(ids))
	for i, id := range ids {
		res[i] = r.Cache.ServiceByID[id]
	}
	return res
}

func (r *queryResolver) allServices() []*cache.Service {
	res := make([]*cache.Service, len(r.Cache.SortedServiceIDs))
	for i, id := range r.Cache.SortedActionIDs {
		res[i] = r.Cache.ServiceByID[id]
	}
	return res
}

func (r *queryResolver) Actions(ctx context.Context, filter *string) ([]*cache.Action, error) {
	if filter != nil {
		return r.actionsWithFilter(*filter), nil
	}
	return r.allActions(), nil
}

func (r *queryResolver) actionsWithFilter(filter string) []*cache.Action {
	var ids []string
	for _, id := range r.Cache.SortedActionIDs {
		if !strings.Contains(id, filter) {
			continue
		}
		ids = append(ids, id)
	}
	res := make([]*cache.Action, len(ids))
	for i, id := range ids {
		res[i] = r.Cache.ActionByID[id]
	}
	return res
}

func (r *queryResolver) allActions() []*cache.Action {
	res := make([]*cache.Action, len(r.Cache.SortedActionIDs))
	for i, id := range r.Cache.SortedActionIDs {
		res[i] = r.Cache.ActionByID[id]
	}
	return res
}

func (r *queryResolver) ResourceTypes(ctx context.Context, filter *string) ([]*cache.ResourceType, error) {
	if filter != nil {
		return r.resourceTypesWithFilter(*filter), nil
	}
	return r.allResourceTypes(), nil
}

func (r *queryResolver) resourceTypesWithFilter(filter string) []*cache.ResourceType {
	var ids []string
	for _, id := range r.Cache.SortedResourceTypeIDs {
		if !strings.Contains(id, filter) {
			continue
		}
		ids = append(ids, id)
	}
	res := make([]*cache.ResourceType, len(ids))
	for i, id := range ids {
		res[i] = r.Cache.ResourceTypeByID[id]
	}
	return res
}

func (r *queryResolver) allResourceTypes() []*cache.ResourceType {
	res := make([]*cache.ResourceType, len(r.Cache.SortedResourceTypeIDs))
	for i, id := range r.Cache.SortedResourceTypeIDs {
		res[i] = r.Cache.ResourceTypeByID[id]
	}
	return res
}

func (r *queryResolver) ConditionKeys(ctx context.Context, filter *string) ([]*cache.ConditionKey, error) {
	if filter != nil {
		return r.conditionKeysWithFilter(*filter), nil
	}
	return r.allConditionKeys(), nil
}

func (r *queryResolver) conditionKeysWithFilter(filter string) []*cache.ConditionKey {
	var ids []string
	for _, id := range r.Cache.SortedConditionKeyIDs {
		if !strings.Contains(id, filter) {
			continue
		}
		ids = append(ids, id)
	}
	res := make([]*cache.ConditionKey, len(ids))
	for i, id := range ids {
		res[i] = r.Cache.ConditionKeyByID[id]
	}
	return res
}

func (r *queryResolver) allConditionKeys() []*cache.ConditionKey {
	res := make([]*cache.ConditionKey, len(r.Cache.SortedConditionKeyIDs))
	for i, id := range r.Cache.SortedConditionKeyIDs {
		res[i] = r.Cache.ConditionKeyByID[id]
	}
	return res
}

func (r *resourceTypeResolver) ConditionKeys(ctx context.Context, obj *cache.ResourceType) ([]*cache.ConditionKey, error) {
	res := make([]*cache.ConditionKey, len(obj.ConditionKeys))
	for i, ck := range obj.ConditionKeys {
		res[i] = r.Cache.ConditionKeyByID[ck]
	}
	return res, nil
}

func (r *resourceTypeResolver) Actions(ctx context.Context, obj *cache.ResourceType) ([]*cache.Action, error) {
	aIDs := r.Cache.ActionIDsByResourceTypeID[obj.ID]
	res := make([]*cache.Action, len(aIDs))
	for i, aID := range aIDs {
		res[i] = r.Cache.ActionByID[aID]
	}
	return res, nil
}

func (r *resourceTypeResolver) Service(ctx context.Context, obj *cache.ResourceType) (*cache.Service, error) {
	return r.Cache.ServiceByID[obj.Service], nil
}

func (r *serviceResolver) Actions(ctx context.Context, obj *cache.Service) ([]*cache.Action, error) {
	res := make([]*cache.Action, len(obj.Actions))
	for i, aID := range obj.Actions {
		res[i] = r.Cache.ActionByID[aID]
	}
	return res, nil
}

func (r *serviceResolver) ResourceTypes(ctx context.Context, obj *cache.Service) ([]*cache.ResourceType, error) {
	res := make([]*cache.ResourceType, len(obj.ResourceTypes))
	for i, rtID := range obj.ResourceTypes {
		res[i] = r.Cache.ResourceTypeByID[rtID]
	}
	return res, nil
}

func (r *serviceResolver) ConditionKeys(ctx context.Context, obj *cache.Service) ([]*cache.ConditionKey, error) {
	res := make([]*cache.ConditionKey, len(obj.ConditionKeys))
	for i, ckID := range obj.ConditionKeys {
		res[i] = r.Cache.ConditionKeyByID[ckID]
	}
	return res, nil
}

// Action returns ActionResolver implementation.
func (r *Resolver) Action() ActionResolver { return &actionResolver{r} }

// ActionResourceType returns ActionResourceTypeResolver implementation.
func (r *Resolver) ActionResourceType() ActionResourceTypeResolver {
	return &actionResourceTypeResolver{r}
}

// ConditionKey returns ConditionKeyResolver implementation.
func (r *Resolver) ConditionKey() ConditionKeyResolver { return &conditionKeyResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// ResourceType returns ResourceTypeResolver implementation.
func (r *Resolver) ResourceType() ResourceTypeResolver { return &resourceTypeResolver{r} }

// Service returns ServiceResolver implementation.
func (r *Resolver) Service() ServiceResolver { return &serviceResolver{r} }

type actionResolver struct{ *Resolver }
type actionResourceTypeResolver struct{ *Resolver }
type conditionKeyResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type resourceTypeResolver struct{ *Resolver }
type serviceResolver struct{ *Resolver }
