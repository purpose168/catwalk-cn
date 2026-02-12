package providers

import (
	"slices"
	"testing"
)

// TestValidDefaultModels 验证所有提供商的默认模型ID是否在其模型列表中存在
func TestValidDefaultModels(t *testing.T) {
	for _, p := range GetAll() {
		t.Run(p.Name, func(t *testing.T) {
			var modelIds []string
			for _, m := range p.Models {
				modelIds = append(modelIds, m.ID)
			}
			if !slices.Contains(modelIds, p.DefaultLargeModelID) {
				t.Errorf("提供商 %q 的默认大模型 %q 未在模型列表中找到", p.Name, p.DefaultLargeModelID)
			}
			if !slices.Contains(modelIds, p.DefaultSmallModelID) {
				t.Errorf("提供商 %q 的默认小模型 %q 未在模型列表中找到", p.Name, p.DefaultSmallModelID)
			}
		})
	}
}
