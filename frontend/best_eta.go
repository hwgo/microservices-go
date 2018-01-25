package frontend

import (
	"context"
	// "errors"
	// "math"
	// "sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hwgo/pher/log"
	"github.com/hwgo/pher/pool"

	"github.com/hwgo/config"
	"github.com/hwgo/driver"
)

type bestETA struct {
	// customer customer.Interface
	driver driver.Interface
	// route    route.Interface
	pool   *pool.Pool
	logger log.Factory
}

// Response contains ETA for a trip.
type Response struct {
	Driver string
	ETA    time.Duration
}

func newBestETA(tracer opentracing.Tracer, logger log.Factory) *bestETA {
	return &bestETA{
		// customer: customer.NewClient(
		// 	tracer,
		// 	logger.With(zap.String("component", "customer_client")),
		// ),
		driver: driver.NewClient(
			tracer,
			logger.With(zap.String("component", "driver_client")),
		),
		// route: route.NewClient(
		// 	tracer,
		// 	logger.With(zap.String("component", "route_client")),
		// ),
		pool:   pool.New(config.RouteWorkerPoolSize),
		logger: logger,
	}
}

func (eta *bestETA) Get(ctx context.Context, customerID string) (*Response, error) {
	eta.logger.For(ctx).Info("bestETA@frontend")

	// customer, err := eta.customer.Get(ctx, customerID)
	// if err != nil {
	// 	return nil, err
	// }
	// eta.logger.For(ctx).Info("Found customer", zap.Any("customer", customer))

	customer := struct {
		ID       string
		Name     string
		Location string
	}{
		ID:       "123",
		Name:     "Rachel's Floral Designs",
		Location: "115,277",
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetBaggageItem("customer", customer.Name)
	}

	drivers, err := eta.driver.FindNearest(ctx, customer.Location)
	if err != nil {
		return nil, err
	}
	eta.logger.For(ctx).Info("Found drivers", zap.Any("drivers", drivers))

	// results := eta.getRoutes(ctx, customer, drivers)
	// eta.logger.For(ctx).Info("Found routes", zap.Any("routes", results))

	// resp := &Response{ETA: math.MaxInt64}
	// for _, result := range results {
	// 	if result.err != nil {
	// 		return nil, err
	// 	}
	// 	if result.route.ETA < resp.ETA {
	// 		resp.ETA = result.route.ETA
	// 		resp.Driver = result.driver
	// 	}
	// }
	// if resp.Driver == "" {
	// 	return nil, errors.New("No routes found")
	// }

	resp := &Response{
		ETA:    time.Duration(3) * time.Minute,
		Driver: "BRU670",
	}

	eta.logger.For(ctx).Info("Dispatch successful", zap.String("driver", resp.Driver), zap.String("eta", resp.ETA.String()))
	return resp, nil
}

// type routeResult struct {
// 	driver string
// 	route  *route.Route
// 	err    error
// }

// getRoutes calls Route service for each (customer, driver) pair
// func (eta *bestETA) getRoutes(ctx context.Context, customer *customer.Customer, drivers []driver.Driver) []routeResult {
// 	results := make([]routeResult, 0, len(drivers))
// 	wg := sync.WaitGroup{}
// 	routesLock := sync.Mutex{}
// 	for _, dd := range drivers {
// 		wg.Add(1)
// 		driver := dd // capture loop var
// 		// Use worker pool to (potentially) execute requests in parallel
// 		eta.pool.Execute(func() {
// 			route, err := eta.route.FindRoute(ctx, driver.Location, customer.Location)
// 			routesLock.Lock()
// 			results = append(results, routeResult{
// 				driver: driver.DriverID,
// 				route:  route,
// 				err:    err,
// 			})
// 			routesLock.Unlock()
// 			wg.Done()
// 		})
// 	}
// 	wg.Wait()
// 	return results
// }
