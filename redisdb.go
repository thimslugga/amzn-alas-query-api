package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis"
)

// Hash key for the full ALAS cache in Redis
const alasCacheKey = "__alas_cache"

// global db object
var db *Database

// Database is the base struct for the Redis db
type Database struct {
	Client *redis.Client
	Ready  bool
}

// NewDatabase returns a new Redis db instance
func NewDatabase() (database *Database, err error) {
	database = &Database{}
	database.Client = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       config.RedisDatabase,
	})

	err = database.Client.Ping().Err()
	return
}

// SetReady sets the db as ready to serve requests
func (db *Database) SetReady() {
	db.Ready = true
}

// RefreshALASLoop is the main loop for keeping the ALAS in sync with Redis.
// On entry the cache is refreshed once, then a timer is started to re-check
// on the interval provided in the environment or the default 300 seconds.
func (db *Database) RefreshALASLoop() {
	var err error
	if err = db.RefreshALASCache(); err != nil {
		log.Fatal("Failed to refresh the ALAS cache, ", err)
	}
	db.SetReady()
	ticker := time.NewTicker(time.Duration(config.CacheTTL) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case _ = <-ticker.C:
			if err = db.RefreshALASCache(); err != nil {
				log.Fatal("Failed to refresh the ALAS cache, ", err)
			}
		}
	}
}

// RefreshALASCache will retrieve the ALAS feed for each amznlinux release
// and update the Redis cache accordingly.
func (db *Database) RefreshALASCache() error {
	log.Println("Starting refresh of ALAS cache")
	vulns := make(map[string]interface{}, 0)
	pkgToVulnMap := make(map[string][]string, 0)
	for _, feed := range []string{amazonLinuxFeed, amazonLinux2Feed} {
		log.Println("Retrieving feed from", feed)
		feed, err := GetALASFeed(feed)
		if err != nil {
			return err
		}
		for _, x := range feed.Channel.Vulns {
			if !db.ALASExists(x.ALASString()) {
				alasString := x.ALASString()
				log.Println("Expanding", alasString)
				expanded := x.Expand()
				vulns[alasString] = expanded.ToJSON()
				for _, pkg := range expanded.Packages {
					if _, ok := pkgToVulnMap[pkg]; !ok {
						pkgToVulnMap[pkg] = make([]string, 0)
					}
					pkgToVulnMap[pkg] = append(pkgToVulnMap[pkg], alasString)
				}
			}
		}
	}
	var err error
	if len(vulns) > 0 {
		log.Println("Writing newly discovered vulnerabilities to db")
		if err = db.AddVulns(vulns, pkgToVulnMap); err != nil {
			return err
		}
	} else {
		log.Println("No new vulnerabilities discovered")
	}
	return nil
}

// AddVulns writes new vulnerability information to the Redis DB.
// A Hash is used for each ALAS and its corresponding metadata, and a Zset
// is used to map every package to a list of the ALAS's associated with it.
func (db *Database) AddVulns(vulns map[string]interface{}, pkgMap map[string][]string) (err error) {
	if err = db.Client.HMSet(alasCacheKey, vulns).Err(); err != nil {
		return
	}
	for k, v := range pkgMap {
		zs := make([]*redis.Z, 0)
		for _, x := range v {
			zs = append(zs, &redis.Z{Member: x, Score: 0})
		}
		if err = db.Client.ZAddNX(k, zs...).Err(); err != nil {
			return err
		}
	}
	return
}

// ALASExists checks if data exists in Redis for a given ALAS.
func (db *Database) ALASExists(alas string) bool {
	data, err := db.Client.HMGet(alasCacheKey, alas).Result()
	if err != nil || data[0] == nil {
		return false
	}
	return true
}

// GetVulnsByPackage retrieves a list of ALAS strings associated with a given
// package.
func (db *Database) GetVulnsByPackage(pkg string) (vulns []string, err error) {
	vulns, err = db.Client.ZRange(pkg, 0, -1).Result()
	return
}

// GetALAS returns the ExpandedVuln for a given ALAS from Redis.
func (db *Database) GetALAS(alas string) (vuln ExpandedVuln, err error) {
	data, err := db.Client.HMGet(alasCacheKey, alas).Result()
	if err != nil || data[0] == nil {
		return
	}
	json.Unmarshal([]byte(data[0].(string)), &vuln)
	return
}
