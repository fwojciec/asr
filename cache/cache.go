package cache

import (
	"fmt"
	"sort"

	"github.com/fwojciec/asr"
)

type Service struct {
	ID            string
	Prefix        string
	Name          string
	ConfigDocURL  string
	APIDocURL     string
	IAMDocURL     string
	Actions       []string
	ResourceTypes []string
	ConditionKeys []string
}

type Action struct {
	ID               string
	Name             string
	DocURL           string
	Description      string
	AccessLevel      string
	ResourceTypes    []ActionResourceType
	ConditionKeys    []string
	DependentActions []string
	Service          string
}

type ResourceType struct {
	ID            string
	Name          string
	DocURL        string
	ARNPattern    string
	ConditionKeys []string
	Service       string
}

type ConditionKey struct {
	ID          string
	Name        string
	DocURL      string
	Description string
	Type        string
}

type ActionResourceType struct {
	ResourceType string
	Required     bool
}

type Cache struct {
	ServiceByID                     map[string]*Service
	ActionByID                      map[string]*Action
	ResourceTypeByID                map[string]*ResourceType
	ConditionKeyByID                map[string]*ConditionKey
	ActionIDsByResourceTypeID       map[string][]string
	ActionIDsByConditionKeyID       map[string][]string
	ActionIDsByDependentActionID    map[string][]string
	ResourceTypeIDsByConditionKeyID map[string][]string
	ServiceIDsByConditionKeyID      map[string][]string
	SortedServiceIDs                []string
	SortedActionIDs                 []string
	SortedResourceTypeIDs           []string
	SortedConditionKeyIDs           []string
}

func (c *Cache) addServiceByID(s *Service) {
	if c.ServiceByID == nil {
		c.ServiceByID = make(map[string]*Service)
	}
	c.ServiceByID[s.ID] = s
	c.SortedServiceIDs = sortedAppend(c.SortedServiceIDs, s.ID)
}

func (c *Cache) addActionByID(a *Action) {
	if c.ActionByID == nil {
		c.ActionByID = make(map[string]*Action)
	}
	c.ActionByID[a.ID] = a
	c.SortedActionIDs = sortedAppend(c.SortedActionIDs, a.ID)
}

func (c *Cache) addResourceTypeByID(rt *ResourceType) {
	if c.ResourceTypeByID == nil {
		c.ResourceTypeByID = map[string]*ResourceType{}
	}
	c.ResourceTypeByID[rt.ID] = rt
	c.SortedResourceTypeIDs = sortedAppend(c.SortedResourceTypeIDs, rt.ID)
}

func (c *Cache) addConditionKeyByID(ck *ConditionKey) {
	if c.ConditionKeyByID == nil {
		c.ConditionKeyByID = make(map[string]*ConditionKey)
	}
	c.ConditionKeyByID[ck.ID] = ck
	c.SortedConditionKeyIDs = sortedAppend(c.SortedConditionKeyIDs, ck.ID)
}

func (c *Cache) addActionIDByResourceTypeID(rtID, aID string) {
	if c.ActionIDsByResourceTypeID == nil {
		c.ActionIDsByResourceTypeID = make(map[string][]string)
	}
	sortedMapAppend(rtID, aID, c.ActionIDsByResourceTypeID)
}

func (c *Cache) addActionIDByConditionKeyID(ckID, aID string) {
	if c.ActionIDsByConditionKeyID == nil {
		c.ActionIDsByConditionKeyID = make(map[string][]string)
	}
	sortedMapAppend(ckID, aID, c.ActionIDsByConditionKeyID)
}

func (c *Cache) addActionIDByDependentActionID(daID, aID string) {
	if c.ActionIDsByDependentActionID == nil {
		c.ActionIDsByDependentActionID = make(map[string][]string)
	}
	sortedMapAppend(daID, aID, c.ActionIDsByDependentActionID)
}

func (c *Cache) addResourceTypeIDsByConditionKeyID(ckID, rtID string) {
	if c.ResourceTypeIDsByConditionKeyID == nil {
		c.ResourceTypeIDsByConditionKeyID = make(map[string][]string)
	}
	sortedMapAppend(ckID, rtID, c.ResourceTypeIDsByConditionKeyID)
}

func (c *Cache) addServiceIDsByConditionKeyID(ckID, sID string) {
	if c.ServiceIDsByConditionKeyID == nil {
		c.ServiceIDsByConditionKeyID = make(map[string][]string)
	}
	sortedMapAppend(ckID, sID, c.ServiceIDsByConditionKeyID)
}

func (c *Cache) AddService(s *asr.Service) {
	service := &Service{
		ID:           s.Prefix,
		Prefix:       s.Prefix,
		Name:         s.Name,
		ConfigDocURL: s.ConfigDocURL,
		APIDocURL:    s.APIDocURL,
		IAMDocURL:    s.IAMDocURL,
	}
	c.addServiceActions(s.Actions, service)
	c.addServiceResourceTypes(s.ResourceTypes, service)
	c.addServiceConditionKeys(s.ConditionKeys, service)
	c.addServiceByID(service)
}

func (c *Cache) addServiceActions(as []asr.Action, service *Service) {
	for _, a := range as {
		action := &Action{
			ID:            fmt.Sprintf("%s:%s", service.Prefix, a.Name),
			Name:          a.Name,
			DocURL:        a.DocURL,
			Description:   a.Description,
			AccessLevel:   a.AccessLevel,
			ConditionKeys: a.ConditionKeys,
			Service:       service.ID,
		}
		for _, ck := range a.ConditionKeys {
			c.addActionIDByConditionKeyID(ck, action.ID)
		}
		action.ResourceTypes = make([]ActionResourceType, len(a.ResourceTypes))
		for i, rt := range a.ResourceTypes {
			rtID := fmt.Sprintf("%s:%s", service.Prefix, rt.Name)
			action.ResourceTypes[i] = ActionResourceType{ResourceType: rtID, Required: rt.Required}
			c.addActionIDByResourceTypeID(rtID, action.ID)
		}
		action.DependentActions = make([]string, len(a.DependentActions))
		for j, da := range a.DependentActions {
			action.DependentActions[j] = da
			c.addActionIDByDependentActionID(da, action.ID)
		}
		c.addActionByID(action)
		service.Actions = sortedAppend(service.Actions, action.ID)
	}
}

func (c *Cache) addServiceConditionKeys(cks []asr.ConditionKey, service *Service) {
	for _, ck := range cks {
		conditionKey := &ConditionKey{
			ID:          ck.Name,
			Name:        ck.Name,
			DocURL:      ck.DocURL,
			Description: ck.Description,
			Type:        ck.Type,
		}
		c.addServiceIDsByConditionKeyID(conditionKey.ID, service.ID)
		c.addConditionKeyByID(conditionKey)
		service.ConditionKeys = sortedAppend(service.ConditionKeys, conditionKey.ID)
	}
}

func (c *Cache) addServiceResourceTypes(rts []asr.ResourceType, service *Service) {
	for _, rt := range rts {
		resourceType := &ResourceType{
			ID:            fmt.Sprintf("%s:%s", service.Prefix, rt.Name),
			Name:          rt.Name,
			DocURL:        rt.DocURL,
			ARNPattern:    rt.ARNPattern,
			ConditionKeys: rt.ConditionKeys,
			Service:       service.ID,
		}
		for _, ck := range rt.ConditionKeys {
			c.addResourceTypeIDsByConditionKeyID(ck, resourceType.ID)
		}
		c.addResourceTypeByID(resourceType)
		service.ResourceTypes = sortedAppend(service.ResourceTypes, resourceType.ID)
	}
}

func sortedMapAppend(k string, v string, m map[string][]string) {
	if cur, ok := m[k]; ok {
		m[k] = sortedAppend(cur, v)
	} else {
		m[k] = []string{v}
	}
}

func sortedAppend(data []string, value string) []string {
	if data == nil {
		return []string{value}
	}
	i := sort.Search(len(data), func(i int) bool { return data[i] >= value })
	if i < len(data) && data[i] == value {
		return data
	} else if i == len(data) {
		return append(data, value)
	}
	data = append(data[:i+1], data[i:]...)
	data[i] = value
	return data
}

func NewCache(data []*asr.Service) *Cache {
	c := &Cache{}
	for _, s := range data {
		c.AddService(s)
	}
	return c
}
