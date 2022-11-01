package telemetryservice

import (
	"net/http"
	"os"

	"github.com/edgenesis/shifu/pkg/telemetryservice/mqtt"
	"k8s.io/klog"
)

var serverListenPort = os.Getenv("SERVER_LISTEN_PORT")

// TODO: need to modify path of mqtt.BindMQTTServicehandler after other servie implement
func New(stop <-chan struct{}) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mqtt.BindMQTTServicehandler)
	err := Start(stop, mux, serverListenPort)
	if err != nil {
		klog.Errorf("Error when telemetryService Running, error: %v", err)
	}
}

func Start(stop <-chan struct{}, mux *http.ServeMux, addr string) error {
	var errChan = make(chan error, 1)
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			klog.Errorf("Error when server running, error: %v", err)
			errChan <- err
		}
	}()

	klog.Infof("Listening at %#v", addr)
	select {
	case err := <-errChan:
		return err
	case <-stop:
		return server.Close()
	}
}