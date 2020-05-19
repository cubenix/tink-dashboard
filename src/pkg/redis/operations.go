package redis

import "encoding/json"

// CacheKey is the string key in redis cache
type CacheKey string

// CacheKeys is a collection of allowed keys
var CacheKeys = struct {
	Templates, Hardwares, Workflows, TemplateNames CacheKey
}{
	Templates:     "templates",
	Hardwares:     "hardwares",
	Workflows:     "workflows",
	TemplateNames: "template-names",
}

// Get retrieves the value for a field in a given key
func (c *Cache) Get(key CacheKey, field string) (string, error) {
	res := c.client.HGet(string(key), field)
	if res.Err() != nil {
		return "", res.Err()
	}
	result, err := res.Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetAll retrieves all fields and their respective values in a given key
func (c *Cache) GetAll(key CacheKey) (map[string]string, error) {
	res := c.client.HGetAll(string(key))
	if res.Err() != nil {
		return nil, res.Err()
	}
	result, err := res.Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Set updates the field with the value in a given key
func (c *Cache) Set(key CacheKey, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	res := c.client.HSet(string(key), field, data)
	return res.Err()
}

// SetAll updates the field with the value in a given key
func (c *Cache) SetAll(key CacheKey, fields map[string]interface{}) error {
	res := c.client.HMSet(string(key), fields)
	return res.Err()
}

// Delete deletes fields from a given key
func (c *Cache) Delete(key CacheKey, fields ...string) error {
	res := c.client.HDel(string(key), fields...)
	return res.Err()
}
