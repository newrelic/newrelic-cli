package install

func install() error {
	// Execute the discovery process.
	d := new(mockDiscoverer)
	manifest, err := d.discover()
	if err != nil {
		return err
	}

	// Retrieve the relevant recipes.
	f := new(yamlRecipeFetcher)
	recipes, err := f.fetch(manifest)
	if err != nil {
		return err
	}

	// Execute the recipe steps.

	return nil
}
