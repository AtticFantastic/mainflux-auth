package services

import "fmt"

// AddResource assigns a resource specified using its type and id, to the
// particular owner.
func AddResource(owner, rType, rId string) {
	cKey := fmt.Sprintf("auth:users:%s:owned", owner)
	cVal := fmt.Sprintf("%s:%s", rType, rId)

	c := cache.Get()
	defer c.Close()

	c.Do("SADD", cKey, cVal)
}
