package controller

/*
Copyright 2017-2018 Crunchy Data Solutions, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	//apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	crv1 "github.com/crunchydata/postgres-operator/apis/cr/v1"
	clusteroperator "github.com/crunchydata/postgres-operator/operator/cluster"
)

// PgreplicaController holds the connections for the controller
type PgreplicaController struct {
	PgreplicaClient    *rest.RESTClient
	PgreplicaScheme    *runtime.Scheme
	PgreplicaClientset *kubernetes.Clientset
	Namespace          string
}

// Run starts an pgreplica resource controller
func (c *PgreplicaController) Run(ctx context.Context) error {
	fmt.Print("Watch Pgreplica objects\n")

	_, err := c.watchPgreplicas(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for Pgreplica resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

// watchPgreplicas is the event loop for pgreplica resources
func (c *PgreplicaController) watchPgreplicas(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.PgreplicaClient,
		crv1.PgreplicaResourcePlural,
		//apiv1.NamespaceAll,
		c.Namespace,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		&crv1.Pgreplica{},

		// resyncPeriod
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,

		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		})

	go controller.Run(ctx.Done())
	return controller, nil
}

// onAdd is called when a pgreplica is added
func (c *PgreplicaController) onAdd(obj interface{}) {
	replica := obj.(*crv1.Pgreplica)
	fmt.Printf("[PgreplicaCONTROLLER] OnAdd %s\n", replica.ObjectMeta.SelfLink)
	if replica.Status.State == crv1.PgreplicaStateProcessed {
		log.Info("pgreplica " + replica.ObjectMeta.Name + " already processed")
		return
	}

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use clusterScheme.Copy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	copyObj, err := c.PgreplicaScheme.Copy(replica)
	if err != nil {
		fmt.Printf("ERROR creating a deep copy of replica object: %v\n", err)
		return
	}

	replicaCopy := copyObj.(*crv1.Pgreplica)
	replicaCopy.Status = crv1.PgreplicaStatus{
		State:   crv1.PgreplicaStateProcessed,
		Message: "Successfully processed Pgreplica by controller",
	}

	err = c.PgreplicaClient.Put().
		Name(replica.ObjectMeta.Name).
		Namespace(replica.ObjectMeta.Namespace).
		Resource(crv1.PgreplicaResourcePlural).
		Body(replicaCopy).
		Do().
		Error()

	if err != nil {
		fmt.Printf("ERROR updating status: %v\n", err)
	} else {
		fmt.Printf("UPDATED status: %#v\n", replicaCopy)
	}

	clusteroperator.ScaleBase(c.PgreplicaClientset, c.PgreplicaClient, replicaCopy, replicaCopy.ObjectMeta.Namespace)

}

// onUpdate is called when a pgreplica is updated
func (c *PgreplicaController) onUpdate(oldObj, newObj interface{}) {
	//oldExample := oldObj.(*crv1.Pgreplica)
	newExample := newObj.(*crv1.Pgreplica)
	log.Info("pgreplica " + newExample.ObjectMeta.Name + " updated")

}

// onDelete is called when a pgreplica is deleted
func (c *PgreplicaController) onDelete(obj interface{}) {
	replica := obj.(*crv1.Pgreplica)
	fmt.Printf("[PgreplicaCONTROLLER] OnDelete %s\n", replica.ObjectMeta.SelfLink)
	//clusteroperator.DeleteClusterBase(c.PgclusterClientset, c.PgclusterClient, cluster, cluster.ObjectMeta.Namespace)
}
