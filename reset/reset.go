package reset

import (
	"github.com/hbagdi/deck/utils"
	"github.com/hbagdi/go-kong/kong"
	"github.com/pkg/errors"
)

// Reset deletes all entities in Kong.
func Reset(state *utils.KongRawState, client *kong.Client) error {

	if state == nil {
		return errors.New("state cannot be empty")
	}
	// Delete routes before services
	for _, r := range state.Routes {
		err := client.Routes.Delete(nil, r.ID)
		if err != nil {
			return err
		}
	}

	for _, s := range state.Services {
		err := client.Services.Delete(nil, s.ID)
		if err != nil {
			return err
		}
	}

	for _, c := range state.Consumers {
		err := client.Consumers.Delete(nil, c.ID)
		if err != nil {
			return err
		}
	}

	// Upstreams also removes Targets
	for _, u := range state.Upstreams {
		err := client.Upstreams.Delete(nil, u.ID)
		if err != nil {
			return err
		}
	}

	// Certificates also removes SNIs
	for _, u := range state.Certificates {
		err := client.Certificates.Delete(nil, u.ID)
		if err != nil {
			return err
		}
	}

	for _, u := range state.CACertificates {
		err := client.CACertificates.Delete(nil, u.ID)
		if err != nil {
			if err.Error() == "Not found" {
				return caCertificates, nil
			}
			return err
		}
	}

	for _, p := range state.Plugins {
		// Delete global plugins explicitly since those will not
		// DELETE ON CASCADE
		if p.Consumer == nil && p.Service == nil &&
			p.Route == nil {
			err := client.Plugins.Delete(nil, p.ID)
			if err != nil {
				return err
			}
		}
	}

	// TODO handle custom entities
	return nil
}
