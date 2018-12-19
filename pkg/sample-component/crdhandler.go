/*
Copyright (C) 2018 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package sample-component

import (
	"fmt"
	"strings"

	sample-componentclientset "github.com/blackducksoftware/synopsys-operator/pkg/sample-component/client/clientset/versioned"
	sample-component_v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/sample-component/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/protoform"
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// HandlerInterface contains the methods that are required
type HandlerInterface interface {
	ObjectCreated(obj interface{})
	ObjectDeleted(obj string)
	ObjectUpdated(objOld, objNew interface{})
}

// Handler will store the configuration that is required to initiantiate the informers callback
type Handler struct {
	config      *protoform.Config
	kubeConfig  *rest.Config
	kubeClient  *kubernetes.Clientset
	sample-componentClient *sample-componentclientset.Clientset
	defaults    *sample-component_v1.SampleComponentSpec
}

// NewHandler will create the handler
func NewHandler(config *protoform.Config, kubeConfig *rest.Config, kubeClient *kubernetes.Clientset, sample-componentClient *sample-componentclientset.Clientset, defaults *sample-component_v1.SampleComponentSpec) *Handler {
	return &Handler{config: config, kubeConfig: kubeConfig, kubeClient: kubeClient, sample-componentClient: sample-componentClient, defaults: defaults}
}

// ObjectCreated will be called for create sample-component events
func (h *Handler) ObjectCreated(obj interface{}) {
	log.Debugf("objectCreated: %+v", obj)
	sample-componentv1, ok := obj.(*sample-component_v1.SampleComponent)
	if !ok {
		log.Error("Unable to cast to SampleComponent object")
		return
	}
	if strings.EqualFold(sample-componentv1.Spec.State, "") {
		// merge with default values
		newSpec := sample-componentv1.Spec
		sample-componentDefaultSpec := h.defaults
		err := mergo.Merge(&newSpec, sample-componentDefaultSpec)
		log.Debugf("merged sample-component details %+v", newSpec)
		if err != nil {
			log.Errorf("unable to merge the sample-component structs for %s due to %+v", sample-componentv1.Name, err)
			//Set spec/state  and status/state to started
			h.updateState("error", "error", fmt.Sprintf("unable to merge the sample-component structs for %s due to %+v", sample-componentv1.Name, err), sample-componentv1)
		} else {
			sample-componentv1.Spec = newSpec
			// update status
			sample-componentv1, err := h.updateState("pending", "creating", "", sample-componentv1)

			if err == nil {
				sample-componentCreator := NewCreater(h.kubeConfig, h.kubeClient, h.sample-componentClient)

				// create sample-component instance
				err = sample-componentCreator.CreateSampleComponent(&sample-componentv1.Spec)

				if err != nil {
					h.updateState("error", "error", fmt.Sprintf("%+v", err), sample-componentv1)
				} else {
					h.updateState("running", "running", "", sample-componentv1)
				}
			}
		}
	}
}

// ObjectDeleted will be called for delete sample-component events
func (h *Handler) ObjectDeleted(name string) {
	log.Debugf("objectDeleted: %+v", name)
	sample-componentCreator := NewCreater(h.kubeConfig, h.kubeClient, h.sample-componentClient)
	sample-componentCreator.DeleteSampleComponent(name)
}

// ObjectUpdated will be called for update sample-component events
func (h *Handler) ObjectUpdated(objOld, objNew interface{}) {
	log.Debugf("objectUpdated: %+v", objNew)
}

func (h *Handler) updateState(specState string, statusState string, errorMessage string, sample-component *sample-component_v1.SampleComponent) (*sample-component_v1.SampleComponent, error) {
	sample-component.Spec.State = specState
	sample-component.Status.State = statusState
	sample-component.Status.ErrorMessage = errorMessage
	sample-component, err := h.updateSampleComponentObject(sample-component)
	if err != nil {
		log.Errorf("couldn't update the state of sample-component object: %s", err.Error())
	}
	return sample-component, err
}

func (h *Handler) updateSampleComponentObject(obj *sample-component_v1.SampleComponent) (*sample-component_v1.SampleComponent, error) {
	return h.sample-componentClient.SynopsysV1().SampleComponents(h.config.Namespace).Update(obj)
}
