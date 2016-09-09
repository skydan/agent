package hostInfo

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rancher/agent/model"
	"github.com/rancher/agent/utilities/constants"
	"github.com/rancher/agent/utilities/utils"
)

type OSCollector struct {
	DataGetter OSInfoGetter
	GOOS       string
	InfoData   model.InfoData
}

type OSInfoGetter interface {
	GetOS(model.InfoData) (map[string]string, error)
	GetDockerVersion(model.InfoData, bool) map[string]string
}

type OSDataGetter struct{}

func (o OSDataGetter) GetDockerVersion(infoData model.InfoData, verbose bool) map[string]string {
	data := map[string]string{}
	versionData := infoData.Version
	version := "unknown"
	if verbose && versionData.Version != "" {
		version = fmt.Sprintf("Docker version %v, build %v", versionData.Version, versionData.GitCommit)
	} else if versionData.Version != "" {
		version = utils.SemverTrunk(versionData.Version, 2)
	}
	data["dockerVersion"] = version

	return data
}

func (o OSCollector) GetData() (map[string]interface{}, error) {
	infoData := o.InfoData
	data := map[string]interface{}{}
	osData, err := o.DataGetter.GetOS(infoData)
	if err != nil {
		return data, errors.Wrap(err, constants.OSGetDataError)
	}

	for key, value := range o.DataGetter.GetDockerVersion(infoData, true) {
		data[key] = value
	}
	for key, value := range osData {
		data[key] = value
	}
	return data, nil
}

func (o OSCollector) GetLabels(prefix string) (map[string]string, error) {
	osData, err := o.DataGetter.GetOS(o.InfoData)
	if err != nil {
		return map[string]string{}, errors.Wrap(err, constants.OSGetDataError)
	}
	labels := map[string]string{
		fmt.Sprintf("%s.%s", prefix, "docker_version"):       o.DataGetter.GetDockerVersion(o.InfoData, false)["dockerVersion"],
		fmt.Sprintf("%s.%s", prefix, "linux_kernel_version"): utils.SemverTrunk(osData["kernelVersion"], 2),
	}
	return labels, nil
}

func (o OSCollector) KeyName() string {
	return "osInfo"
}
