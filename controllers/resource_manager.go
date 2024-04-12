package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

type resourcesFromAssets []byte

type Resouces struct {
	ServiceAccount     corev1.ServiceAccount
	Role               rbacv1.Role
	ClusterRole        rbacv1.ClusterRole
	RoleBinding        rbacv1.RoleBinding
	ClusterRoleBinding rbacv1.ClusterRoleBinding
	ConfigMaps         []corev1.ConfigMap
	Daemonset          appsv1.DaemonSet
}

func filePathWalkDir(path string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Path:%s error in filePathWalk: %v", path, err)
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func getResources(path string) []resourcesFromAssets {
	manifests := []resourcesFromAssets{}
	files, err := filePathWalkDir(path)
	if err != nil {
		panic(err)
	}
	sort.Strings(files)
	for _, file := range files {
		buffer, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		manifests = append(manifests, buffer)
	}
	return manifests
}

func addRescourcesControls(path string) (Resouces, controlFunc) {
	res := Resouces{}
	ctrl := controlFunc{}
	fmt.Println("Get assets from: ", path)

	manifests := getResources(path)

	s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme,
		scheme.Scheme, json.SerializerOptions{Yaml: true, Pretty: false, Strict: false})
	reg := regexp.MustCompile(`\b(\w*kind:\w*)\B.*\b`)

	for _, m := range manifests {
		kind := reg.FindString(string(m))
		tmpSlices := strings.Split(kind, ":")
		kind = strings.TrimSpace(tmpSlices[1])

		fmt.Println("Looking for", "Kind", kind, "in path:", path)
		switch kind {
		case "ServiceAccount":
			_, _, err := s.Decode(m, nil, &res.ServiceAccount)
			if err != nil {
				panic(err)
			}
			ctrl = append(ctrl, ServiceAccount)
		case "Role":
			_, _, err := s.Decode(m, nil, &res.Role)
			if err != nil {
				panic(err)
			}
			ctrl = append(ctrl, Role)
		case "ClusterRole":
			_, _, err := s.Decode(m, nil, &res.ClusterRole)
			if err != nil {
				panic(err)
			}
			ctrl = append(ctrl, ClusterRole)
		case "RoleBinding":
			_, _, err := s.Decode(m, nil, &res.RoleBinding)
			if err != nil {
				panic(err)
			}
			ctrl = append(ctrl, RoleBinding)
		case "ClusterRoleBinding":
			_, _, err := s.Decode(m, nil, &res.ClusterRoleBinding)
			if err != nil {
				panic(err)
			}
			ctrl = append(ctrl, ClusterRoleBinding)
		case "ConfigMap":
			cm := corev1.ConfigMap{}
			_, _, err := s.Decode(m, nil, &cm)
			if err != nil {
				panic(err)
			}
			res.ConfigMaps = append(res.ConfigMaps, cm)
			if len(res.ConfigMaps) == 1 {
				ctrl = append(ctrl, ConfigMaps)
			}
		case "DaemonSet":
			_, _, err := s.Decode(m, nil, &res.Daemonset)
			if err != nil {
				panic(err)
			}
			ctrl = append(ctrl, DaemonSet)
		default:
			fmt.Println("Unknown Resource", "Manifest", m, "Kind", kind)
		}
	}

	return res, ctrl
}
