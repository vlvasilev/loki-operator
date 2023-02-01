package manifests

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"

	"github.com/vlvasilev/loki-operator/pkg/manifests/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// buildValiConfigMap creates the single configmap containing the vali and its curator configuration
func buildValiConfigMap(opt *Options) (*corev1.ConfigMap, error) {
	valiConfigBytes, err := buildValiConfiguration(opt.Vali)
	if err != nil {
		return nil, err
	}

	curatorConfigBytes, err := buildCuratorConfiguration(opt.Curator)
	if err != nil {
		return nil, err
	}

	opt.ConfigSHA, err = getConfigSHA(valiConfigBytes, curatorConfigBytes)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      lokiConfigMapName(opt.ConfigSHA),
			Namespace: opt.Namespace,
			Labels:    commonLabels(),
		},
		BinaryData: map[string][]byte{
			"vali.yaml":    valiConfigBytes,
			"curator.yaml": curatorConfigBytes,
		},
	}, nil
}

func getConfigSHA(valiConfigBytes, curatorConfigBytes []byte) (string, error) {
	s := sha1.New()
	_, err := s.Write(valiConfigBytes)
	if err != nil {
		return "", err
	}

	_, err = s.Write(curatorConfigBytes)
	if err != nil {
		return "", err
	}

	sha := fmt.Sprintf("%x", s.Sum(nil))
	return sha[:8], nil
}

// buildValiConfiguration builds a Vali configuration files
func buildValiConfiguration(opts ValiOptions) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	err := config.ValiConfigYAMLTmpl.Execute(w, *opts.Spec)
	if err != nil {
		return nil, fmt.Errorf("failed to create vali configuration: %w", err)
	}
	return ioutil.ReadAll(w)
}

// buildCuratorConfiguration builds a Curator configuration files
func buildCuratorConfiguration(opts CuratorOptions) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	err := config.CuratorConfigYAMLTmpl.Execute(w, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create curator configuration: %w", err)
	}
	return ioutil.ReadAll(w)
}
