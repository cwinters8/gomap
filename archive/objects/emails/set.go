package emails

import (
	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/requests"
)

func (e Email) Set(c *client.Client, acctID string) (requests.Set, error) {
	return requests.NewSet(acctID, e)
	// if err != nil {
	// 	return fmt.Errorf("failed to instantiate new set: %w", err)
	// }
	// req := requests.NewRequest([]requests.Call{s})
	// body, err := req.Send(c)
	// if err != nil {
	// 	return fmt.Errorf("failed to send request: %w", err)
	// }
	// var result Result
	// if err := json.Unmarshal(body, &result); err != nil {
	// 	return fmt.Errorf("failed to unmarshal response body: %w", err)
	// }
	// return nil
}
