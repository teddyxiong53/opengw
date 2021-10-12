/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-13 15:53:31
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 10:19:29
@FilePath: /goAdapter-Raw/device/formula.go
*/
package device

import (
	"fmt"
	"goAdapter/httpServer/model"
	"regexp"
	"strconv"

	"github.com/walkmiao/chili/environment"
)

type Parser interface {
	VarSet(float64) error
	PreVarSet([]model.DeviceTSLPropertyTemplate) error
	SetFormula(formula string)
}

type IndexParser struct {
	formula string
	env     *environment.Environment
}

var _ Parser = (*IndexParser)(nil)

func (fla *IndexParser) SetFormula(formula string) {
	fla.formula = formula
}

func (fla *IndexParser) PreVarSet(variables []model.DeviceTSLPropertyTemplate) error {
	reg, err := regexp.Compile("i([0-9])") //val*i6*i7 i:index
	if err != nil {
		return err
	}
	result := reg.FindAllString(fla.formula, -1)
	for _, item := range result {
		index := item[1:]
		i, err := strconv.Atoi(index)
		if err != nil {
			return err
		}
		nodeVar := variables[i]

		if len(nodeVar.Value) > 0 {
			env.SetIntVariable(item, int64(nodeVar.Value[0].Value.(uint32)))

		} else {
			return fmt.Errorf("基础值%s还未获取,values:%v", item, nodeVar.Value)
		}

	}
	return nil
}

func (fla *IndexParser) VarSet(val float64) error {

	if err := fla.env.SetFloatVariable("val", val); err != nil {
		return fmt.Errorf("设置表达式val值错误:%v", err)
	}
	return nil
}
