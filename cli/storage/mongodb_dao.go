package storage

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBDao struct {
	Session           *mgo.Session
	Idx               mgo.Index
	defaultCollection string
	*SearchDB
}

func NewMongoDBDao(connection string, searchDir string) *MongoDBDao {
	var session *mgo.Session
	var err error

	if trimmed := strings.TrimSuffix(connection, "?ssl=true"); len(trimmed) != len(connection) {
		dialInfo, err := mgo.ParseURL(trimmed)
		if err != nil {
			log.Fatalf("Error trying to parse MongoDB connection string %q: %v", connection, err)
		}

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			tlsConfig := &tls.Config{InsecureSkipVerify: true}
			return tls.Dial("tcp", addr.String(), tlsConfig)
		}

		if session, err = mgo.DialWithInfo(dialInfo); err != nil {
			log.Fatalf("Error trying to connect to MongoDB: %s", err.Error())
		}
	} else if session, err = mgo.Dial(connection); err != nil {
		log.Fatalf("Error trying to connect to MongoDB: %s", err.Error())
	}

	index := mgo.Index{
		Key:    []string{"page_id"},
		Unique: true,
	}

	return &MongoDBDao{
		Session:           session,
		Idx:               index,
		defaultCollection: "root",
		SearchDB:          NewSearchDB(searchDir),
	}
}

func (mongo *MongoDBDao) getCollection(name string) (*mgo.Collection, error) {
	collection := mongo.Session.Copy().DB("").C(name)
	if err := collection.EnsureIndex(mongo.Idx); err != nil {
		return nil, fmt.Errorf("storage.getCollection: failed to ensure index: %v", err)
	}

	return collection, nil
}

func (mongo *MongoDBDao) Fetch(id string) ([]byte, error) {
	if filepath.Ext(id) == "" {
		id = filepath.Join(id, "index.html")
	}

	collection, err := mongo.getCollection(parseCollection(id))
	if err != nil {
		return nil, fmt.Errorf("storage.Fetch: %v", err)
	}

	var p page
	if err := collection.Find(bson.M{"page_id": id}).One(&p); err != nil {
		return nil, fmt.Errorf("storage.Fetch: %v", err)
	}

	if filepath.Ext(p.ID) == ".html" {
		return []byte(html.UnescapeString(p.Data)), nil
	}
	return []byte(p.Data), nil
}

func (mongo *MongoDBDao) FetchGlob(pattern string) []string {
	var results []string

	collection, err := mongo.getCollection(parseCollection(pattern))
	if err != nil {
		return results
	}

	pages := make([]page, 0)
	if err := collection.Find(nil).Select(bson.M{"id": 0, "page_id": 1}).All(&pages); err != nil {
		return results
	}

	for _, p := range pages {
		results = append(results, p.ID)
	}

	return results
}

func (mongo *MongoDBDao) Insert(id string, data []byte) error {
	collection, err := mongo.getCollection(parseCollection(id))
	if err != nil {
		return fmt.Errorf("storage.Insert: %v", err)
	}

	if _, err := collection.Upsert(bson.M{"page_id": id}, bson.M{"$set": NewPage(id, data)}); err != nil {
		return fmt.Errorf("storage.Insert: %v", err)
	}

	return nil
}

func (mongo *MongoDBDao) Delete(id string) error {
	collection, err := mongo.getCollection(parseCollection(id))
	if err != nil {
		return fmt.Errorf("storage.Delete: %v", err)
	}

	if err := mongo.SearchDB.Index.Delete(id); err != nil {
		return fmt.Errorf("storage.Delete: %v", err)
	}

	if err := collection.Remove(bson.M{"page_id": id}); err != nil {
		return fmt.Errorf("storage.Delete: %v", err)
	}

	return nil
}

func (mongo *MongoDBDao) DeleteGlob(pattern string) error {
	collection, err := mongo.getCollection(parseCollection(pattern))
	if err != nil {
		return fmt.Errorf("storage.DeleteGlob: %v", err)
	}

	var pages []page
	if err := collection.Find(bson.M{"route": collection.Name}).All(&pages); err != nil {
		return fmt.Errorf("storage.Delete: %v", err)
	}

	for _, p := range pages {
		if err := mongo.SearchDB.Index.Delete(p.ID); err != nil {
			return fmt.Errorf("storage.Delete: %v", err)
		}
	}

	if err := collection.DropCollection(); err != nil {
		return fmt.Errorf("storage.Delete: %v", err)
	}

	return nil
}

func (mongo *MongoDBDao) Index(id, title string, data []byte) error {
	globalData.Lock()
	defer globalData.Unlock()

	indexData := struct {
		Title    string    `json:"title"`
		Body     string    `json:"body"`
		Modified time.Time `json:"modified"`
	}{
		Title:    title,
		Body:     string(filterHTMLTags(data)),
		Modified: time.Now(),
	}

	if err := mongo.SearchDB.Index.Index(id, indexData); err != nil {
		return fmt.Errorf("storage.Index: %v", err)
	}

	return nil
}

func (mongo *MongoDBDao) Query(query string) (*QueryResult, error) {
	globalData.RLock()
	defer globalData.RUnlock()
	return mongo.SearchDB.Query(query)
}

type page struct {
	ID    string `json:"page_id"`
	Route string `json:"page_route"`
	Data  string `json:"page_data"`
}

func NewPage(id string, data []byte) page {
	p := page{ID: id, Route: parseCollection(id)}
	if filepath.Ext(id) == ".html" {
		p.Data = html.EscapeString(string(json.RawMessage(data)))
	} else {
		p.Data = string(data)
	}
	return p
}

func parseCollection(pageID string) string {
	paths := strings.Split(pageID, "/")
	if paths[0] == "" && len(paths) > 1 {
		if filepath.Ext(paths[1]) != "" {
			return "root"
		}
		return paths[1]
	}

	return paths[0]
}
