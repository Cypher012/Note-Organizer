package automate

import "log"

func UserFlow() error {
	client := CreateClient()
	user := NewTestUser(false)

	// if err := Register(client, user); err != nil {
	// 	return err
	// }

	if err := Login(client, user); err != nil {
		return err
	}

	if err := GetFolders(client); err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		if err := CreateFolder(client); err != nil {
			return err
		}
	}

	if err := GetFolders(client); err != nil {
		return err
	}

	log.Println("All passed")
	return nil
}
