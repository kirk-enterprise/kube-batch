package kirkutil

import (
	"strconv"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
)

const (
	CPUOvercommitRatioAnnotation = "k8s.qiniu.com/cpu-overcommit-ratio"
	CPUOvercommitRatioMin        = 0.1
	CPUOvercommitRatioMax        = 10
)

// GetCPUOvercommitRatio returns CPU over-commmit ratio of node.
func GetCPUOvercommitRatio(node *v1.Node) float64 {
	if ratio, ok := node.Annotations[CPUOvercommitRatioAnnotation]; ok {
		ratio_f, err := strconv.ParseFloat(ratio, 64)
		if err != nil {
			glog.Errorf("failed to parse ratio data: %s", ratio)
			return 1.0
		}
		if ratio_f < CPUOvercommitRatioMin {
			ratio_f = CPUOvercommitRatioMin
		} else if ratio_f > CPUOvercommitRatioMax {
			ratio_f = CPUOvercommitRatioMax
		}
		return ratio_f
	}
	return 1.0
}

func UpdateAllocatable(rl v1.ResourceList, node *v1.Node) v1.ResourceList {
	newrl := make(v1.ResourceList)
	ratio := GetCPUOvercommitRatio(node)
	for rName, rQuant := range rl {
		switch rName {
		case v1.ResourceCPU:
			rQuant.SetMilli(int64(float64(rQuant.MilliValue()) * ratio))
		}
		newrl[rName] = rQuant
	}
	return newrl
}
