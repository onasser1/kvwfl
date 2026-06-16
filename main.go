package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	logger "github.com/onasser1/kontroller/pkg/apis/okofs/v1alpha1"
	"github.com/onasser1/validating-kontroller/pkg/options"
	"github.com/spf13/pflag"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/klog/v2"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func main() {
	opts := options.NewDefaultOptions()

	flagSet := pflag.NewFlagSet(options.ValidatingKontroller, pflag.ExitOnError)
	globalflag.AddGlobalFlags(flagSet, options.ValidatingKontroller)
	opts.AddFlagSet(flagSet)

	if err := flagSet.Parse(os.Args); err != nil {
		klog.Fatalf("error reading flags from flagset: %v", err)
	}
	c := opts.Config()

	// mux is the http handler we would use for the https server
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(ServeValidationAdmission))

	ch := server.SetupSignalHandler()
	stoppedCh, listerStoppedCh, err := c.SecureServingInfo.Serve(mux, 30*time.Second, ch)
	if err != nil {
		klog.Fatalf("error from Serve function: %v", err)
	} else {
		<-stoppedCh
		<-listerStoppedCh
	}

}

func ServeValidationAdmission(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hi, we are ServeValidationAdmission function, we are called")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		klog.Fatalf("error reading body from request: %v", err)
	}
	// Get GVK of AdmissionReview
	admissionReviewGVK := admissionv1beta1.SchemeGroupVersion.WithKind("AdmissionReview")
	AdmissionReview := admissionv1beta1.AdmissionReview{}
	_, _, err = codecs.UniversalDeserializer().Decode(body, &admissionReviewGVK, &AdmissionReview)
	if err != nil {
		klog.Fatalf("error deserializing the AdmissionReview object: %v", err)
	}

	// Get GVK of Logger type
	loggerGVK := logger.SchemaGroupVersion.WithKind("Logger")
	Logger := logger.Logger{}
	_, _, err = codecs.UniversalDeserializer().Decode(AdmissionReview.Request.Object.Raw, &loggerGVK, &Logger)
	if err != nil {
		klog.Fatalf("error deserializing the Logger object: %v", err)
	}
	klog.Infof("Logger Resource we got is: %v", Logger)

	AdmissionResponse := admissionv1beta1.AdmissionResponse{}
	if allowed := ValidateLogger(&Logger); !allowed {
		AdmissionResponse = admissionv1beta1.AdmissionResponse{
			UID:     AdmissionReview.Request.UID,
			Allowed: allowed,
			Result: &v1.Status{
				Message: "Not Allowed",
			},
		}
	} else {
		AdmissionResponse = admissionv1beta1.AdmissionResponse{
			UID:     AdmissionReview.Request.UID,
			Allowed: allowed,
		}
	}
	klog.Infof("AdmissionResponse Allowed? : %v", AdmissionResponse.Allowed)
	AdmissionReview.Response = &AdmissionResponse

	response, err := json.Marshal(AdmissionReview)
	if err != nil {
		klog.Fatalf("error marshaling AdmissionReview back to JSON: %v", err)
	}
	if _, err = w.Write(response); err != nil {
		klog.Fatalf("error writing AdmissionReview reponse: %v", err)
	}
}

// ValidateLogger currently returns boolean result
// we can return a status message with the reason
// to be used within AdmissionReview.Result
func ValidateLogger(l *logger.Logger) bool {
	labels := l.GetLabels()
	if labels == nil {
		return true
	}
	if labels["system"] == "kube" || labels["env"] == "production" {
		return false
	}
	return true
}
