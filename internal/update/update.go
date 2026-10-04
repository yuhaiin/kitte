package update

import (
	"context"
	"fmt"
	"log"
)

func Run(ctx context.Context, root string) error {
	client := newHTTPClient()

	log.Printf("download geoip database")
	geoIP, err := client.get(ctx, geoIPURL)
	if err != nil {
		return err
	}
	if err := generateGeoIP(root, geoIP); err != nil {
		return err
	}

	log.Printf("download geosite database")
	geoSite, err := client.get(ctx, geoSiteURL)
	if err != nil {
		return err
	}
	if err := generateGeoSite(root, geoSite); err != nil {
		return err
	}

	downloads := make(map[string][]byte)
	for _, source := range ruleSources() {
		log.Printf("download %s", source.name)
		data, err := client.get(ctx, source.url)
		if err != nil {
			return fmt.Errorf("%s: %w", source.name, err)
		}
		downloads[source.url] = data
	}
	if err := updateRules(root, downloads); err != nil {
		return err
	}

	log.Printf("download Tailscale DERP map")
	derp, err := client.get(ctx, tailscaleURL)
	if err != nil {
		return err
	}
	if err := generateTailscale(root, derp); err != nil {
		return err
	}

	log.Printf("generate manifest")
	if err := generateManifest(root); err != nil {
		return err
	}

	log.Printf("update complete")
	return nil
}
