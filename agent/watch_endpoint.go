package agent

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/consul/api/watch"
	"github.com/hashicorp/go-uuid"
)

// Watch ...
type Watch struct {
	Watch map[string]interface{} `json:"watches,omitempty" hcl:"watches" mapstructure:"watches"`
}

// WatchConfiguration apply watch
func (s *HTTPServer) WatchConfiguration(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	args := Watch{}

	if req.Method == "POST" || req.Method == "PUT" {
		if err := decodeBody(req.Body, &args.Watch); err != nil {
			return nil, BadRequestError{err.Error()}
		}
	}

	switch req.Method {
	case "POST":
		return s.agent.addWatch(&resp, &args)
	case "PUT":
		return s.agent.updateWatch(&resp, &args, req.URL.Query().Get("id"))
	case "DELETE":
		return s.agent.deleteWatch(&resp, req.URL.Query().Get("id"))
	default:
		return nil, MethodNotAllowedError{req.Method, []string{"POST", "PUT", "DELETE"}}
	}
}

func (a *Agent) addWatch(resp *http.ResponseWriter, watchconfig *Watch) (interface{}, error) {

	// Compile the watches
	if handlerType, ok := watchconfig.Watch["handler_type"]; !ok {
		watchconfig.Watch["handler_type"] = "http"
	} else if handlerType != "http" {
		return nil, BadRequestError{"Only http handler type supported"}
	}

	// Don't let people use connect watches via this mechanism for now as it
	// needs thought about how to do securely and shouldn't be necessary. Note
	// that if the type assertion fails an type is not a string then
	// ParseExample below will error so we don't need to handle that case.
	if typ, ok := watchconfig.Watch["type"].(string); ok {
		if strings.HasPrefix(typ, "connect_") {
			return nil, BadRequestError{fmt.Sprintf("Watch type %s is not allowed in agent config", typ)}
		}
	}

	// Parse the watches, excluding 'handler' and 'args'
	wp, err := watch.ParseExempt(watchconfig.Watch, []string{"handler", "args"})
	if err != nil {
		return nil, BadRequestError{fmt.Sprintf("Failed to parse watch (%#v): %v", watchconfig.Watch, err)}
	}

	config, err := a.config.APIConfig(true)
	if err != nil {
		a.logger.Printf("[ERR] agent: Failed to run watch: %v", err)
		return nil, BadRequestError{fmt.Sprintf("[ERR] agent: Failed to run watch: %v", err)}
	}

	wp.UUID, _ = uuid.GenerateUUID()
	// Store the watch plan
	// Fire off a goroutine for each new watch plan.
	a.watchPlans = append(a.watchPlans, wp)
	go func(wp *watch.Plan) {
		httpConfig := wp.Exempt["http_handler_config"].(*watch.HttpHandlerConfig)
		wp.Handler = makeHTTPWatchHandler(a.LogOutput, httpConfig)
		wp.LogOutput = a.LogOutput
		addr := config.Address
		if config.Scheme == "https" {
			addr = "https://" + addr
		}
		if err := wp.RunWithConfig(addr, config); err != nil {
			a.logger.Printf("[ERR] agent: Failed to run watch: %v", err)
		}
	}(wp)

	return wp.UUID, nil
}

func (a *Agent) updateWatch(resp *http.ResponseWriter, watchconfig *Watch, id string) (interface{}, error) {
	index := -1
	watchplans := a.watchPlans
	for i, wp := range watchplans {
		if wp.UUID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return nil, BadRequestError{"watch not found"}
	}

	// Compile the watches
	if handlerType, ok := watchconfig.Watch["handler_type"]; !ok {
		watchconfig.Watch["handler_type"] = "http"
	} else if handlerType != "http" {
		return nil, BadRequestError{fmt.Sprintf("Only http handler type supported")}
	}

	// Don't let people use connect watches via this mechanism for now as it
	// needs thought about how to do securely and shouldn't be necessary. Note
	// that if the type assertion fails an type is not a string then
	// ParseExample below will error so we don't need to handle that case.
	if typ, ok := watchconfig.Watch["type"].(string); ok {
		if strings.HasPrefix(typ, "connect_") {
			return nil, BadRequestError{fmt.Sprintf("Watch type %s is not allowed in agent config", typ)}
		}
	}

	// Parse the watches, excluding 'handler' and 'args'
	wp, err := watch.ParseExempt(watchconfig.Watch, []string{"handler", "args"})
	if err != nil {
		return nil, BadRequestError{fmt.Sprintf("Failed to parse watch (%#v): %v", watchconfig.Watch, err)}
	}

	config, err := a.config.APIConfig(true)
	if err != nil {
		a.logger.Printf("[ERR] agent: Failed to run watch: %v", err)
		return nil, BadRequestError{fmt.Sprintf("[ERR] agent: Failed to run watch: %v", err)}
	}

	wp.UUID, _ = uuid.GenerateUUID()
	// Store the watch plan
	// Fire off a goroutine for each new watch plan.
	watchplans[index] = wp
	go func(wp *watch.Plan) {
		httpConfig := wp.Exempt["http_handler_config"].(*watch.HttpHandlerConfig)
		wp.Handler = makeHTTPWatchHandler(a.LogOutput, httpConfig)
		wp.LogOutput = a.LogOutput
		addr := config.Address
		if config.Scheme == "https" {
			addr = "https://" + addr
		}
		if err := wp.RunWithConfig(addr, config); err != nil {
			a.logger.Printf("[ERR] agent: Failed to run watch: %v", err)
		}
	}(wp)

	return wp.UUID, nil
}

func (a *Agent) deleteWatch(resp *http.ResponseWriter, watchID string) (interface{}, error) {

	index := -1
	for i, wp := range a.watchPlans {
		if wp.UUID == watchID {
			index = i
			break
		}
	}

	if index == -1 {
		return nil, BadRequestError{"watch not found"}
	}

	a.watchPlans[index] = a.watchPlans[len(a.watchPlans)-1]
	a.watchPlans[len(a.watchPlans)-1] = nil
	a.watchPlans = a.watchPlans[:len(a.watchPlans)-1]

	return nil, nil
}
