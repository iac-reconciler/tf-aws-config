package compare

import (
	"github.com/iac-reconciler/tf-aws-config/pkg/load"
	log "github.com/sirupsen/logrus"
)

type LocatedItem struct {
	load.ConfigurationItem
	Location string `json:"location"`
}

// Reconcile reconcile the snapshot and tfstates.
// Not yet implemented, so returns an empty struct
func Reconcile(snapshot load.Snapshot, tfstates map[string]load.TerraformState) (items []LocatedItem, err error) {
	// first load each item into memory, referenced by resourceType, resourceId, ARN
	var (
		configTypeIdMap = make(map[string]map[string]*load.ConfigurationItem)
		configArnMap    = make(map[string]*load.ConfigurationItem)
		// the keys are resource types, using the AWS-Config keys;
		// the values are map[string]string
		// in there, the keys are arn or id (if no arn), the values are location,
		// one of "terraform", "config", "both"
		itemToLocation      = make(map[string]map[string]string)
		filteredConfigItems []*load.ConfigurationItem
	)
	for _, item := range snapshot.ConfigurationItems {
		if item.ResourceType == configComplianceResourceType {
			continue
		}
		if item.ResourceType == "" {
			log.Warnf("AWS Config snapshot: empty resource type for item %s", item.ARN)
			continue
		}
		if _, ok := configTypeIdMap[item.ResourceType]; !ok {
			configTypeIdMap[item.ResourceType] = make(map[string]*load.ConfigurationItem)
		}
		if item.ResourceID != "" {
			configTypeIdMap[item.ResourceType][item.ResourceID] = &item
		}
		if item.ARN != "" {
			configArnMap[item.ARN] = &item
		}
		filteredConfigItems = append(filteredConfigItems, &item)
	}
	// now comes the harder part. We have to go through each tfstate and reconcile it with the snapshot
	// This would be easy if there were standards, but everything is driven by the provider,
	// terraform itself has no standard or intelligence about it, so we need to know all of them.
	for statefile, tfstate := range tfstates {
		for i, resource := range tfstate.Resources {
			// only care about managed resources
			if resource.Mode != load.TerraformManaged {
				continue
			}
			// only care about aws resources
			if resource.Provider != "provider.aws" && resource.Provider != `provider["registry.terraform.io/hashicorp/aws"]"` {
				continue
			}
			// look up the resource type
			var (
				configType string
				ok         bool
			)
			if configType, ok = awsTerraformToConfigTypeMap[resource.Type]; !ok {
				log.Warnf("unknown terraform resource type: %s", resource.Type)
				continue
			}
			if _, ok := itemToLocation[configType]; !ok {
				itemToLocation[configType] = make(map[string]string)
			}
			if _, ok := configTypeIdMap[configType]; !ok {
				configTypeIdMap[configType] = make(map[string]*load.ConfigurationItem)
			}
			for j, instance := range resource.Instances {
				var (
					resourceId, arn, location string
					item                      *load.ConfigurationItem
				)
				// try by arn first
				arnPtr := instance.Attributes["arn"]
				if arnPtr != nil {
					arn = arnPtr.(string)
				}
				resourceIdPtr := instance.Attributes["id"]
				if resourceIdPtr != nil {
					resourceId = resourceIdPtr.(string)
				}
				if arn != "" {
					item = configArnMap[arn]
				}
				if item == nil && resourceId != "" {
					item = configTypeIdMap[configType][resourceId]
				}

				// if we could not find it by ARN or by configType+id, then
				// it is only in terraform
				if item == nil {
					location = "terraform"
					item = &load.ConfigurationItem{
						ResourceType: configType,
						ResourceID:   resourceId,
						ARN:          arn,
					}
					// if it only is in Terraform, then we do not have it in the maps
					// so we need to add it for later reference
					if arn != "" {
						configArnMap[arn] = item
					}
					if resourceId != "" {
						configTypeIdMap[configType][resourceId] = item
					}
				} else {
					location = "both"
				}

				if location == "" {
					// no resource ID or arn, which shouldn't occur, so warn
					log.Warnf("unable to find resource ID or ARN for resource %d, instance %d in file %s", i, j, statefile)
				} else {
					key := arn
					if key == "" {
						key = resourceId
					}
					itemToLocation[configType][key] = location
				}
			}
		}
	}
	// go through config and see what is not covered already in terraform
	for _, item := range filteredConfigItems {
		if _, ok := itemToLocation[item.ResourceType]; !ok {
			itemToLocation[item.ResourceType] = make(map[string]string)
		}
		key := item.ARN
		if key == "" {
			key = item.ResourceID
		}
		if _, ok := itemToLocation[item.ResourceType][key]; !ok {
			itemToLocation[item.ResourceType][key] = "config"
		}
	}
	// now we have all of the resources listed as in Terraform, Config or both
	// so create the reconciled info
	for resourceType, locations := range itemToLocation {
		for key, location := range locations {
			var configItem *load.ConfigurationItem
			if item, ok := configArnMap[key]; ok {
				configItem = item
			} else if item, ok := configTypeIdMap[resourceType][key]; ok {
				configItem = item
			} else {
				log.Errorf("could not find item %s in config maps", key)
				continue
			}
			locatedItem := LocatedItem{
				ConfigurationItem: *configItem,
				Location:          location,
			}
			items = append(items, locatedItem)
		}
	}
	return items, nil
}
