package app

func RunSSHRouter() {
	go func() {
		if err := runSSHRouter(); err != nil {
			panic(err)
		}
	}()
}

func runSSHRouter() error {
	return nil
}
