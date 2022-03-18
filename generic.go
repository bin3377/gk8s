package gk8s

import (
	"context"
	"fmt"
	"io"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/kubectl/pkg/scheme"

	"k8s.io/client-go/kubernetes"
)

type Resource interface {
	runtime.Object
	metav1.Object
}

func Get[T Resource](apiClient kubernetes.Interface, namespace, name string, opt metav1.GetOptions) (T, error) {
	result := new(T)
	switch any(result).(type) {

	case **corev1.Namespace:
		item, err := apiClient.CoreV1().Namespaces().Get(context.Background(), name, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(item))
		return *result, nil

	case **corev1.Pod:
		item, err := apiClient.CoreV1().Pods(namespace).Get(context.Background(), name, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(item))
		return *result, nil

	case **corev1.Secret:
		item, err := apiClient.CoreV1().Secrets(namespace).Get(context.Background(), name, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(item))
		return *result, nil

	case **appsv1.Deployment:
		item, err := apiClient.AppsV1().Deployments(namespace).Get(context.Background(), name, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(item))
		return *result, nil

	case **batchv1.Job:
		item, err := apiClient.BatchV1().Jobs(namespace).Get(context.Background(), name, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(item))
		return *result, nil

	case **rbacv1.Role:
		item, err := apiClient.RbacV1().Roles(namespace).Get(context.Background(), name, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(item))
		return *result, nil

	case **rbacv1.RoleBinding:
		item, err := apiClient.RbacV1().RoleBindings(namespace).Get(context.Background(), name, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(item))
		return *result, nil

	default:
		return *result, fmt.Errorf("%T is not supported", result)

	}
}

func List[T Resource](apiClient kubernetes.Interface, namespace string, opt metav1.ListOptions) ([]T, error) {
	result := make([]T, 0)

	switch any(result).(type) {
	case []*corev1.Namespace:
		list, err := apiClient.CoreV1().Namespaces().List(context.Background(), opt)
		if err != nil {
			return nil, err
		}
		result := make([]T, len(list.Items))
		for i := range list.Items {
			reflect.ValueOf(result).Index(i).Set(reflect.ValueOf(&list.Items[i]))
		}
		return result, nil

	case []*corev1.Pod:
		list, err := apiClient.CoreV1().Pods(namespace).List(context.Background(), opt)
		if err != nil {
			return nil, err
		}
		result := make([]T, len(list.Items))
		for i := range list.Items {
			reflect.ValueOf(result).Index(i).Set(reflect.ValueOf(&list.Items[i]))
		}
		return result, nil

	case []*corev1.Secret:
		list, err := apiClient.CoreV1().Secrets(namespace).List(context.Background(), opt)
		if err != nil {
			return nil, err
		}
		result := make([]T, len(list.Items))
		for i := range list.Items {
			reflect.ValueOf(result).Index(i).Set(reflect.ValueOf(&list.Items[i]))
		}
		return result, nil

	case []*appsv1.Deployment:
		list, err := apiClient.AppsV1().Deployments(namespace).List(context.Background(), opt)
		if err != nil {
			return nil, err
		}
		result := make([]T, len(list.Items))
		for i := range list.Items {
			reflect.ValueOf(result).Index(i).Set(reflect.ValueOf(&list.Items[i]))
		}
		return result, nil

	case []*batchv1.Job:
		list, err := apiClient.BatchV1().Jobs(namespace).List(context.Background(), opt)
		if err != nil {
			return nil, err
		}
		result := make([]T, len(list.Items))
		for i := range list.Items {
			reflect.ValueOf(result).Index(i).Set(reflect.ValueOf(&list.Items[i]))
		}
		return result, nil

	case []*rbacv1.Role:
		list, err := apiClient.RbacV1().Roles(namespace).List(context.Background(), opt)
		if err != nil {
			return nil, err
		}
		result := make([]T, len(list.Items))
		for i := range list.Items {
			reflect.ValueOf(result).Index(i).Set(reflect.ValueOf(&list.Items[i]))
		}
		return result, nil

	case []*rbacv1.RoleBinding:
		list, err := apiClient.RbacV1().RoleBindings(namespace).List(context.Background(), opt)
		if err != nil {
			return nil, err
		}
		result := make([]T, len(list.Items))
		for i := range list.Items {
			reflect.ValueOf(result).Index(i).Set(reflect.ValueOf(&list.Items[i]))
		}
		return result, nil

	default:
		return nil, fmt.Errorf("%T is not supported", result)
	}

}

func Create[T Resource](apiClient kubernetes.Interface, namespace string, spec T, opt metav1.CreateOptions) (T, error) {
	result := new(T)

	switch any(result).(type) {
	case **corev1.Namespace:
		var obj corev1.Namespace
		reflect.ValueOf(&obj).Elem().Set(reflect.ValueOf(spec).Elem())
		newObj, err := apiClient.CoreV1().Namespaces().Create(context.Background(), &obj, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(newObj))
		return *result, nil

	case **corev1.Pod:
		var obj corev1.Pod
		reflect.ValueOf(&obj).Elem().Set(reflect.ValueOf(spec).Elem())
		newObj, err := apiClient.CoreV1().Pods(namespace).Create(context.Background(), &obj, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(newObj))
		return *result, nil

	case **corev1.Secret:
		var obj corev1.Secret
		reflect.ValueOf(&obj).Elem().Set(reflect.ValueOf(spec).Elem())
		newObj, err := apiClient.CoreV1().Secrets(namespace).Create(context.Background(), &obj, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(newObj))
		return *result, nil

	case **appsv1.Deployment:
		var obj appsv1.Deployment
		reflect.ValueOf(&obj).Elem().Set(reflect.ValueOf(spec).Elem())
		newObj, err := apiClient.AppsV1().Deployments(namespace).Create(context.Background(), &obj, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(newObj))
		return *result, nil

	case **batchv1.Job:
		var obj batchv1.Job
		reflect.ValueOf(&obj).Elem().Set(reflect.ValueOf(spec).Elem())
		newObj, err := apiClient.BatchV1().Jobs(namespace).Create(context.Background(), &obj, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(newObj))
		return *result, nil

	case **rbacv1.Role:
		var obj rbacv1.Role
		reflect.ValueOf(&obj).Elem().Set(reflect.ValueOf(spec).Elem())
		newObj, err := apiClient.RbacV1().Roles(namespace).Create(context.Background(), &obj, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(newObj))
		return *result, nil

	case **rbacv1.RoleBinding:
		var obj rbacv1.RoleBinding
		reflect.ValueOf(&obj).Elem().Set(reflect.ValueOf(spec).Elem())
		newObj, err := apiClient.RbacV1().RoleBindings(namespace).Create(context.Background(), &obj, opt)
		if err != nil {
			return *result, err
		}
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(newObj))
		return *result, nil

	default:
		return *result, fmt.Errorf("%T is not supported", result)
	}

}

func Delete[T Resource](apiClient kubernetes.Interface, namespace, name string, opt metav1.DeleteOptions) error {
	hint := new(T)

	switch any(hint).(type) {
	case **corev1.Namespace:
		return apiClient.CoreV1().Namespaces().Delete(context.Background(), name, opt)

	case **corev1.Pod:
		return apiClient.CoreV1().Pods(namespace).Delete(context.Background(), name, opt)

	case **corev1.Secret:
		return apiClient.CoreV1().Secrets(namespace).Delete(context.Background(), name, opt)

	case **appsv1.Deployment:
		return apiClient.AppsV1().Deployments(namespace).Delete(context.Background(), name, opt)

	case **batchv1.Job:
		return apiClient.BatchV1().Jobs(namespace).Delete(context.Background(), name, opt)

	case **rbacv1.Role:
		return apiClient.RbacV1().Roles(namespace).Delete(context.Background(), name, opt)

	case **rbacv1.RoleBinding:
		return apiClient.RbacV1().RoleBindings(namespace).Delete(context.Background(), name, opt)

	default:
		return fmt.Errorf("%T is not supported", hint)
	}

}

func DecodeIntoString[T Resource](input string, into T) error {
	return DecodeInto([]byte(input), into)
}

func DecodeInto[T Resource](data []byte, into T) error {
	decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDecoder()
	return runtime.DecodeInto(decoder, data, into)
}

func PrintYAML[T Resource](obj T, w io.Writer) error {
	return new(printers.YAMLPrinter).PrintObj(obj, w)
}

func PrintJSON[T Resource](obj T, w io.Writer) error {
	return new(printers.JSONPrinter).PrintObj(obj, w)
}
