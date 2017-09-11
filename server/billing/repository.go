package billing

import (
	"context"

	"log"

	"encoding/json"

	"github.com/coreos/etcd/client"
)

func getProject(name string) *Project {
	p, err := Api.Get(context.Background(), "projects/"+name, nil)

	if err != nil {
		if client.IsKeyNotFound(err) {
			return nil
		}

		log.Fatal("Error reading project from etcd. ", err.Error())
	}

	var project Project
	err := json.Unmarshal([]byte(p.Node.Value), &project)
	if err != nil {
		log.Fatal("Error decoding json from etcd: ", err.Error())
	}

	return &project
}

func saveProject(project Project) {
	json, err := json.Marshal(project)
	if err != nil {
		log.Fatal("Error encoding json for etcd: ", err.Error())
	}
	_, err = Api.Set(context.Background(), "projects/" + project.Name, string(json), nil)
	if err != nil {
		log.Fatal("Error saving to etcd: ", err.Error())
	}
}