/*
Copyright 2022 The KCP Authors.

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

package cache

import (
	"context"
	"time"

	"github.com/kcp-dev/logicalcluster/v2"
	"k8s.io/client-go/tools/cache"
)

// ScopeableSharedIndexInformer is an informer that knows how to scope itself down to one cluster,
// or act as an informer across clusters.
type ScopeableSharedIndexInformer interface {
	Cluster(cluster logicalcluster.Name) cache.SharedIndexInformer
	cache.SharedIndexInformer
}

type CancellableEventHandlerRegistrar interface {
	// AddCancellableEventHandler adds an event handler to the shared informer using the shared informer's resync
	// period.  Events to a single handler are delivered sequentially, but there is no coordination
	// between different handlers.
	AddCancellableEventHandler(handler CancellableResourceEventHandler)
	// AddCancellableEventHandlerWithResyncPeriod adds an event handler to the
	// shared informer with the requested resync period; zero means
	// this handler does not care about resyncs.  The resync operation
	// consists of delivering to the handler an update notification
	// for every object in the informer's local cache; it does not add
	// any interactions with the authoritative storage.  Some
	// informers do no resyncs at all, not even for handlers added
	// with a non-zero resyncPeriod.  For an informer that does
	// resyncs, and for each handler that requests resyncs, that
	// informer develops a nominal resync period that is no shorter
	// than the requested period but may be longer.  The actual time
	// between any two resyncs may be longer than the nominal period
	// because the implementation takes time to do work and there may
	// be competing load and scheduling noise.
	AddCancellableEventHandlerWithResyncPeriod(handler CancellableResourceEventHandler, resyncPeriod time.Duration)
}

// CancellableResourceEventHandler can handle notifications for events that
// happen to a resource. The events are informational only, so you
// can't return an error.  The handlers MUST NOT modify the objects
// received; this concerns not only the top level of structure but all
// the data structures reachable from it.
//  * OnAdd is called when an object is added.
//  * OnUpdate is called when an object is modified. Note that oldObj is the
//      last known state of the object-- it is possible that several changes
//      were combined together, so you can't use this to see every single
//      change. OnUpdate is also called when a re-list happens, and it will
//      get called even if nothing changed. This is useful for periodically
//      evaluating or syncing something.
//  * OnDelete will get the final state of the item if it is known, otherwise
//      it will get an object of type DeletedFinalStateUnknown. This can
//      happen if the watch is closed and misses the delete event and we don't
//      notice the deletion until the subsequent re-list.
//  * Done returns a channel which, when closed, indicates that the handler
//      is not expecting to handle any future events
type CancellableResourceEventHandler interface {
	OnAdd(obj interface{})
	OnUpdate(oldObj, newObj interface{})
	OnDelete(obj interface{})
	Done() <-chan struct{}
}

type cancellableResourceEventHandler struct {
	cache.ResourceEventHandler
	ctx context.Context
}

func (c *cancellableResourceEventHandler) Done() <-chan struct{} {
	return c.ctx.Done()
}

// WithContext adds cancellation to a resource handler.
func WithContext(ctx context.Context, handler cache.ResourceEventHandler) CancellableResourceEventHandler {
	return &cancellableResourceEventHandler{
		ResourceEventHandler: handler,
		ctx:                  ctx,
	}
}
