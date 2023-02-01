package manifests

import (
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BuildAll builds all manifests required to run a Loki Stack
func BuildAll(opts *Options) ([]client.Object, error) {
	res := make([]client.Object, 0)

	cm, err := buildValiConfigMap(opts)
	if err != nil {
		return nil, err
	}

	monolithValiSts, err := buildMonolithStatefulSet(opts)
	if err != nil {
		return nil, err
	}

	monolithValiService := buildMonolithValiService(opts)

	res = append(res, cm)
	res = append(res, monolithValiSts)
	res = append(res, monolithValiService)

	return res, nil
}

func buildMonolithStatefulSet(opts *Options) (*appsv1.StatefulSet, error) {
	// Get the Vali stuffs
	valiContainer := getValiContainer(&opts.Vali)
	valiVolumes := getValiVolumes(opts)
	valiPVC := getValiVolumeClaimTemplate(&opts.Vali)
	// Get the curator stuff
	curatorContainer := getCuratorContainer(&opts.Curator)

	//TODO: (vlvasilev) make merge of all containers, volumes and PVCs from all the components in future.

	sts, err := getEmptyValiStatefulSet(opts)
	if err != nil {
		return nil, err
	}

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, valiContainer)
	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, curatorContainer)

	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, valiVolumes...)

	sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, valiPVC)

	return sts, nil
}
