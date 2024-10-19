/*
This file is part of the Tensegrity distribution (https://github.com/fastforgeinc/tensegrity)
Copyright (C) 2024 FastForge Inc.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// KeysNotProducedReason is added in Tensegrity resource when keys are not fully produced.
	KeysNotProducedReason = "KeysNotProduced"
	// KeysNotProducedMessage is added in Tensegrity resource when keys are not produced.
	KeysNotProducedMessage = "Keys are not produced: %s."
	// KeysProducedReason is added in Tensegrity resource when keys are fully produced.
	KeysProducedReason = "KeysProduced"
	// KeysProducedMessage is added in Tensegrity resource when keys are fully produced.
	KeysProducedMessage = "All keys are produced."
	// KeysNotConsumedReason is added in Tensegrity resource when keys are not fully consumed.
	KeysNotConsumedReason = "KeysNotConsumed"
	// KeysNotConsumedMessage is added in Tensegrity resource when keys are not consumed.
	KeysNotConsumedMessage = "Keys are not consumed for envs: %s."
	// KeysConsumedReason is added in Tensegrity resource when keys are fully consumed.
	KeysConsumedReason = "KeyConsumed"
	// KeysConsumedMessage is added in Tensegrity resource when keys are fully consumed.
	KeysConsumedMessage = "All keys are consumed."
)

// TensegrityConditionType defines the conditions of Tensegrity resource.
type TensegrityConditionType string

const (
	// TensegrityConsumed means keys are fully consumed and values are found.
	TensegrityConsumed TensegrityConditionType = "Consumed"
	// TensegrityProduced means keys are fully produced and values are found.
	TensegrityProduced TensegrityConditionType = "Produced"
)

type TensegrityCondition struct {
	// Type of Tensegrity resource condition.
	Type TensegrityConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// LastUpdateTime is the last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
	// LastTransitionTime is a time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// Reason for the condition's last transition.
	Reason string `json:"reason"`
	// Message is a human-readable message indicating details about the transition.
	Message string `json:"message"`
}

func NewTensegrityCondition(
	typ TensegrityConditionType, status corev1.ConditionStatus, reason, message string) *TensegrityCondition {

	return &TensegrityCondition{
		Type:               typ,
		Status:             status,
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

func SetTensegrityCondition(status *TensegrityStatus, condition TensegrityCondition) bool {
	currentCond := GetTensegrityCondition(*status, condition.Type)
	if currentCond != nil && currentCond.Status == condition.Status &&
		currentCond.Reason == condition.Reason && currentCond.Message == condition.Message {
		return false
	}
	if currentCond != nil && currentCond.Status == condition.Status {
		condition.LastTransitionTime = currentCond.LastTransitionTime
	}
	newConditions := filterOutCondition(status.Conditions, condition.Type)
	status.Conditions = append(newConditions, condition)
	return true
}

func GetTensegrityCondition(status TensegrityStatus, typ TensegrityConditionType) *TensegrityCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == typ {
			return &c
		}
	}
	return nil
}

func RemoveTensegrityCondition(status *TensegrityStatus, typ TensegrityConditionType) {
	status.Conditions = filterOutCondition(status.Conditions, typ)
}

func filterOutCondition(conditions []TensegrityCondition, typ TensegrityConditionType) []TensegrityCondition {
	var newConditions []TensegrityCondition
	for _, c := range conditions {
		if c.Type == typ {
			continue
		}
		newConditions = append(newConditions, c)
	}
	return newConditions
}
