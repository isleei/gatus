package api

import endpointui "github.com/TwiN/gatus/v5/config/endpoint/ui"

type ManagedEndpointUIConfigPayload struct {
	HideConditions              bool `json:"hideConditions,omitempty"`
	HideHostname                bool `json:"hideHostname,omitempty"`
	HideURL                     bool `json:"hideURL,omitempty"`
	HidePort                    bool `json:"hidePort,omitempty"`
	HideErrors                  bool `json:"hideErrors,omitempty"`
	DontResolveFailedConditions bool `json:"dontResolveFailedConditions,omitempty"`
	ResolveSuccessfulConditions bool `json:"resolveSuccessfulConditions,omitempty"`
}

func managedEndpointUIConfigToPayload(config *endpointui.Config) *ManagedEndpointUIConfigPayload {
	if config == nil {
		return nil
	}
	return &ManagedEndpointUIConfigPayload{
		HideConditions:              config.HideConditions,
		HideHostname:                config.HideHostname,
		HideURL:                     config.HideURL,
		HidePort:                    config.HidePort,
		HideErrors:                  config.HideErrors,
		DontResolveFailedConditions: config.DontResolveFailedConditions,
		ResolveSuccessfulConditions: config.ResolveSuccessfulConditions,
	}
}

func managedEndpointUIConfigFromPayload(payload *ManagedEndpointUIConfigPayload) *endpointui.Config {
	if payload == nil {
		return nil
	}
	return &endpointui.Config{
		HideConditions:              payload.HideConditions,
		HideHostname:                payload.HideHostname,
		HideURL:                     payload.HideURL,
		HidePort:                    payload.HidePort,
		HideErrors:                  payload.HideErrors,
		DontResolveFailedConditions: payload.DontResolveFailedConditions,
		ResolveSuccessfulConditions: payload.ResolveSuccessfulConditions,
	}
}
