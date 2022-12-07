package deep

import "github.com/revel/revel"

type Depth struct {
	*revel.Controller
}

func (c Depth) DoDeep() revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"deep": 1,
	})
}