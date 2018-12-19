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

package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SampleComponent is a specification for a SampleComponent resource
type SampleComponent struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SampleComponentSpec   `json:"spec"`
	Status SampleComponentStatus `json:"status,omitempty"`
}

// SampleComponentSpec is the spec for a SampleComponent resource
type SampleComponentSpec struct {
	Namespace string `json:"namespace"`
	State     string `json:"state"`

	EnableMetrics bool `json:"enableMetrics"`

	// CPU and memory configurations
	// Example: "300m"
	DefaultCPU string `json:"defaultCpu,omitempty"`
	// Example: "1300Mi"
	DefaultMem string `json:"defaultMem,omitempty"`

	// Log level
	LogLevel string `json:"logLevel,omitempty"`

	ConfigMapName string `json:"configMapName"`

	// Configuration secret
	SecretName string `json:"secretName"`
}

// SampleComponentStatus is the status for a SampleComponent resource
type SampleComponentStatus struct {
	State        string `json:"state"`
	ErrorMessage string `json:"errorMessage"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SampleComponentList is a list of SampleComponent resources
type SampleComponentList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`

	Items []SampleComponent `json:"items"`
}
