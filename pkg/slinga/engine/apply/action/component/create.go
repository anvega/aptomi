package component

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/slinga/engine/apply/action"
	"github.com/Aptomi/aptomi/pkg/slinga/eventlog"
	"github.com/Aptomi/aptomi/pkg/slinga/object"
	"time"
)

type CreateAction struct {
	object.Metadata
	*action.Base

	ComponentKey string
}

func NewCreateAction(componentKey string) *CreateAction {
	return &CreateAction{
		Metadata:     object.Metadata{}, // TODO: initialize
		Base:         action.NewBase(),
		ComponentKey: componentKey,
	}
}

func (a *CreateAction) Apply(context *action.Context) error {
	// deploy to cloud
	err := a.processDeployment(context)
	if err != nil {
		context.EventLog.LogError(err)
		return fmt.Errorf("Errors while creating component '%s': %s", a.ComponentKey, err)
	}

	// update actual state
	return a.updateActualState(context)
}

func (a *CreateAction) updateActualState(context *action.Context) error {
	// get instance from desired state
	instance := context.DesiredState.ComponentInstanceMap[a.ComponentKey]

	// update creation and update times
	instance.UpdateTimes(time.Now(), time.Now())

	// copy it over to the actual state
	context.ActualState.ComponentInstanceMap[a.ComponentKey] = instance
	err := context.ActualStateUpdater.Create(instance)
	if err != nil {
		return fmt.Errorf("error while update actual state: %s", err)
	}
	return nil
}

func (a *CreateAction) processDeployment(context *action.Context) error {
	instance := context.DesiredState.ComponentInstanceMap[a.ComponentKey]
	component := context.DesiredPolicy.Services[instance.Key.ServiceName].GetComponentsMap()[instance.Key.ComponentName]

	if component == nil {
		// This is a service instance. Do nothing
		return nil
	}

	// Instantiate component
	context.EventLog.WithFields(eventlog.Fields{
		"componentKey": instance.Key,
		"component":    component.Name,
		"code":         instance.CalculatedCodeParams,
	}).Info("Deploying new component instance: " + instance.Key.GetKey())

	if component.Code != nil {
		clusterName, ok := instance.CalculatedCodeParams["cluster"].(string)
		if !ok {
			return fmt.Errorf("No cluster specified in code params, component instance: %v", a.ComponentKey)
		}

		cluster, ok := context.DesiredPolicy.Clusters[clusterName]
		if !ok {
			return fmt.Errorf("Can't find cluster in policy: %s", clusterName)
		}

		plugin, err := context.Plugins.GetDeployPlugin(component.Code.Type)
		if err != nil {
			return err
		}

		err = plugin.Create(cluster, a.ComponentKey, instance.CalculatedCodeParams, context.EventLog)
		if err != nil {
			return err
		}
	}

	return nil
}
