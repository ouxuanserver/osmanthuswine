package core

import "encoding/json"

type Controller struct {
	Response
	Request Request
}

func (c *Controller) ControllerInit(req Request, res Response) {
	c.Request = req
	c.Session = res.Session
	c.OriginResponseWriter = res.OriginResponseWriter
}

func (c *Controller) RequestToStruct(v interface{}) error {
	if c.Request.BODY != "" {
		err := json.Unmarshal([]byte(c.Request.BODY), v)
		if err == nil {
			return err
		}
	}
	data, _ := json.Marshal(c.Request.REQUEST)
	return json.Unmarshal(data, v)
}
