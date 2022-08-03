/*
Copyright 2021 The Knative Authors

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

package sinkbinding

import (
	"k8s.io/apimachinery/pkg/util/uuid"
	"knative.dev/eventing/test/rekt/resources/deployment"
	"knative.dev/eventing/test/rekt/resources/pingsource"
	"knative.dev/eventing/test/rekt/resources/sinkbinding"
	"knative.dev/reconciler-test/pkg/eventshub"
	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/manifest"
	"knative.dev/reconciler-test/resources/svc"
)

func SinkBindingV1Deployment() *feature.Feature {
	sbinding := feature.MakeRandomK8sName("sinkbinding")
	sink := feature.MakeRandomK8sName("sink")
	subject := feature.MakeRandomK8sName("subject")
	extensionSecret := string(uuid.NewUUID())

	f := feature.NewFeatureNamed("SinkBinding goes ready")

	f.Setup("install sink", eventshub.Install(sink, eventshub.StartReceiver))
	f.Setup("install a service", svc.Install(sink, "app", "rekt"))
	f.Setup("install a deployment", deployment.Install(subject))

	extensions := map[string]string{
		"sinkbinding": extensionSecret,
	}

	cfg := []manifest.CfgFn{
		sinkbinding.WithExtensions(extensions),
		pingsource.WithSink(svc.AsKReference(sink), ""),
	}

	f.Setup("install SinkBinding", sinkbinding.Install(sbinding, svc.AsDestinationRef(sink), deployment.AsTrackerReference(subject), cfg...))
	f.Setup("SinkBinding goes ready", sinkbinding.IsReady(sbinding))

	return f
}
