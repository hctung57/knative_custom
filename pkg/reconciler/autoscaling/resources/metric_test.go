/*
Copyright 2019 The Knative Authors

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

package resources

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/serving/pkg/apis/autoscaling"
	autoscalingv1alpha1 "knative.dev/serving/pkg/apis/autoscaling/v1alpha1"
	asconfig "knative.dev/serving/pkg/autoscaler/config"
	"knative.dev/serving/pkg/autoscaler/config/autoscalerconfig"
	. "knative.dev/serving/pkg/testing"
)

func TestMakeMetric(t *testing.T) {
	cases := []struct {
		name string
		pa   *autoscalingv1alpha1.PodAutoscaler
		msn  string
		want *autoscalingv1alpha1.Metric
	}{{
		name: "defaults",
		pa:   pa(),
		msn:  "ik",
		want: metric(withScrapeTarget("ik")),
	}, {
		name: "with too short panic window",
		pa:   pa(WithWindowAnnotation("10s"), WithPanicWindowPercentageAnnotation("10")),
		msn:  "wil",
		want: metric(withScrapeTarget("wil"), withWindowAnnotation("10s"),
			withStableWindow(10*time.Second), withPanicWindow(asconfig.BucketSize),
			withPanicWindowPercentageAnnotation("10")),
	}, {
		name: "with longer stable window, no panic window percentage, defaults to 10%",
		pa:   pa(WithWindowAnnotation("10m")),
		msn:  "nu",
		want: metric(
			withScrapeTarget("nu"),
			withStableWindow(10*time.Minute), withPanicWindow(time.Minute),
			withWindowAnnotation("10m")),
	}, {
		name: "with longer panic window percentage",
		pa:   pa(WithPanicWindowPercentageAnnotation("50")),
		msn:  "dansen",
		want: metric(
			withScrapeTarget("dansen"),
			withStableWindow(time.Minute), withPanicWindow(30*time.Second),
			withPanicWindowPercentageAnnotation("50")),
	}, {
		name: "with panic window percentage+rounding",
		pa:   pa(WithPanicWindowPercentageAnnotation("51")),
		msn:  "dansen",
		want: metric(
			withScrapeTarget("dansen"),
			withStableWindow(time.Minute), withPanicWindow(31*time.Second),
			withPanicWindowPercentageAnnotation("51")),
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if diff := cmp.Diff(tc.want, MakeMetric(tc.pa, tc.msn, config)); diff != "" {
				t.Errorf("%q (-want, +got):\n%v", tc.name, diff)
			}
		})
	}
}

func TestStableWindow(t *testing.T) {
	// Not set on PA.
	thePa := pa()
	if got, want := StableWindow(thePa, config), config.StableWindow; got != want {
		t.Errorf("StableWindow = %v, want: %v", got, want)
	}

	thePa = pa(WithWindowAnnotation("251s"))
	if got, want := StableWindow(thePa, config), 251*time.Second; got != want {
		t.Errorf("StableWindow = %v, want: %v", got, want)
	}
}

type MetricOption func(*autoscalingv1alpha1.Metric)

func withStableWindow(window time.Duration) MetricOption {
	return func(metric *autoscalingv1alpha1.Metric) {
		metric.Spec.StableWindow = window
	}
}

func withPanicWindow(window time.Duration) MetricOption {
	return func(metric *autoscalingv1alpha1.Metric) {
		metric.Spec.PanicWindow = window
	}
}

func withWindowAnnotation(window string) MetricOption {
	return func(metric *autoscalingv1alpha1.Metric) {
		metric.Annotations[autoscaling.WindowAnnotationKey] = window
	}
}

func withPanicWindowPercentageAnnotation(percentage string) MetricOption {
	return func(metric *autoscalingv1alpha1.Metric) {
		metric.Annotations[autoscaling.PanicWindowPercentageAnnotationKey] = percentage
	}
}

func withScrapeTarget(s string) MetricOption {
	return func(metric *autoscalingv1alpha1.Metric) {
		metric.Spec.ScrapeTarget = s
	}
}

func metric(options ...MetricOption) *autoscalingv1alpha1.Metric {
	m := &autoscalingv1alpha1.Metric{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-namespace",
			Name:      "test-name",
			Annotations: map[string]string{
				autoscaling.ClassAnnotationKey: autoscaling.KPA,
			},
			Labels:          map[string]string{},
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(pa())},
		},
		Spec: autoscalingv1alpha1.MetricSpec{
			StableWindow: 60 * time.Second,
			PanicWindow:  6 * time.Second,
		},
	}
	for _, fn := range options {
		fn(m)
	}
	return m
}

func pa(options ...PodAutoscalerOption) *autoscalingv1alpha1.PodAutoscaler {
	p := &autoscalingv1alpha1.PodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-namespace",
			Name:      "test-name",
			Annotations: map[string]string{
				autoscaling.ClassAnnotationKey: autoscaling.KPA,
			},
		},
		Spec: autoscalingv1alpha1.PodAutoscalerSpec{
			ContainerConcurrency: 0,
		},
		Status: autoscalingv1alpha1.PodAutoscalerStatus{},
	}
	for _, fn := range options {
		fn(p)
	}
	return p
}

var config = &autoscalerconfig.Config{
	EnableScaleToZero:                  true,
	ContainerConcurrencyTargetFraction: 1.0,
	ContainerConcurrencyTargetDefault:  100.0,
	RPSTargetDefault:                   200.0,
	TargetUtilization:                  0.7,
	MaxScaleUpRate:                     10.0,
	StableWindow:                       60 * time.Second,
	PanicThresholdPercentage:           200,
	PanicWindowPercentage:              10,
	ScaleToZeroGracePeriod:             30 * time.Second,
}
