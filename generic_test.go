package gk8s

import (
	_ "embed"
	"strings"
	"time"

	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...any) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d: "+msg+"\n\n", append([]any{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// OK fails the test if an err is not nil.
func OK(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d: unexpected error: %s\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act any) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

var (
	testConf, _ = FromKubeConfig()
	apiClient   = testConf.APIClient
	//go:embed testdata/pod.yaml
	podYAML []byte
	//go:embed testdata/pod.json
	podJSON []byte
)

func Test_CRUD(t *testing.T) {
	var spec corev1.Pod
	err := DecodeInto(podYAML, &spec)
	OK(t, err)
	spec.Name = fmt.Sprintf("busybox-sleep-%d", time.Now().Unix())
	pNew, err := Create(apiClient, "default", &spec, metav1.CreateOptions{})
	OK(t, err)
	t.Cleanup(func() { Delete[*corev1.Pod](apiClient, "default", pNew.Name, *metav1.NewDeleteOptions(0)) })
	Equals(t, spec.Name, pNew.Name)
	ps, err := List[*corev1.Pod](apiClient, "default", metav1.ListOptions{})
	OK(t, err)
	Assert(t, len(ps) > 1, "get pods")
	pGet, err := Get[*corev1.Pod](apiClient, "default", spec.Name, metav1.GetOptions{})
	OK(t, err)
	Equals(t, spec.Name, pGet.Name)
	err = Delete[*corev1.Pod](apiClient, "default", pNew.Name, metav1.DeleteOptions{})
	OK(t, err)
}

func Test_YAML(t *testing.T) {
	var p corev1.Pod
	err := DecodeInto(podYAML, &p)
	OK(t, err)
	Equals(t, "busybox-sleep", p.Name)
	w := new(strings.Builder)
	err = PrintYAML(&p, w)
	OK(t, err)
	// fmt.Println(w.String())
}

func Test_JSON(t *testing.T) {
	var p corev1.Pod
	err := DecodeInto(podJSON, &p)
	OK(t, err)
	Equals(t, "busybox-sleep", p.Name)
	w := new(strings.Builder)
	err = PrintJSON(&p, w)
	OK(t, err)
	// fmt.Println(w.String())
}
