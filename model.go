package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

// connect to redis
func getDB() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client, nil
}

// SetValue func set specific value for specific key
func SetValue(key string, val int) error {
	client, err := getDB()
	err = client.Set(key, val, 0).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// GEtValue func get value of specific key
func GetValue(key string) (string, error) {
	client, err := getDB()
	val, err := client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
