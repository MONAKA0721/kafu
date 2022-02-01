package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var revision = "undefined"

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic '%v' captured\n", err)
		}
	}()

	fmt.Printf("Version is %s\n", revision)

	// Aura requires you to use "neo4j+s" protocol
	// (You need to replace your connection details, username and password)
	uri := os.Getenv("NEO4J_URI")
	auth := neo4j.BasicAuth(os.Getenv("NEO4J_USERNAME"), os.Getenv("NEO4J_PASSWORD"), "")
	// You typically have one driver instance for the entire application. The
	// driver maintains a pool of database connections to be used by the sessions.
	// The driver is thread safe.
	driver, err := neo4j.NewDriver(uri, auth)
	if err != nil {
		panic(err)
	}
	// Don't forget to close the driver connection when you are finished with it
	defer driver.Close()
	// Create a session to run transactions in. Sessions are lightweight to
	// create and use. Sessions are NOT thread safe.
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	// // WriteTransaction retries the operation in case of transient errors by
	// // invoking the anonymous function multiple times until it succeeds.
	// records, err := session.WriteTransaction(
	// 	func(tx neo4j.Transaction) (interface{}, error) {
	// 		// To learn more about the Cypher syntax, see https://neo4j.com/docs/cypher-manual/current/
	// 		// The Reference Card is also a good resource for keywords https://neo4j.com/docs/cypher-refcard/current/
	// 		createRelationshipBetweenPeopleQuery := `
	// 			MERGE (p1:Person { name: $person1_name })
	// 			MERGE (p2:Person { name: $person2_name })
	// 			MERGE (p1)-[:KNOWS]->(p2)
	// 			RETURN p1, p2`
	// 		result, err := tx.Run(createRelationshipBetweenPeopleQuery, map[string]interface{}{
	// 			"person1_name": "Alice",
	// 			"person2_name": "David",
	// 		})
	// 		if err != nil {
	// 			// Return the error received from driver here to indicate rollback,
	// 			// the error is analyzed by the driver to determine if it should try again.
	// 			return nil, err
	// 		}
	// 		// Collects all records and commits the transaction (as long as
	// 		// Collect doesn't return an error).
	// 		// Beware that Collect will buffer the records in memory.
	// 		return result.Collect()
	// 	})
	// if err != nil {
	// 	panic(err)
	// }
	// for _, record := range records.([]*neo4j.Record) {
	// 	firstPerson := record.Values[0].(neo4j.Node)
	// 	fmt.Printf("First: '%s'\n", firstPerson.Props["name"].(string))
	// 	secondPerson := record.Values[1].(neo4j.Node)
	// 	fmt.Printf("Second: '%s'\n", secondPerson.Props["name"].(string))
	// }

	// Now read the created persons. By using ReadTransaction method a connection
	// to a read replica can be used which reduces load on writer nodes in cluster.
	_, err = session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Code within this function might be invoked more than once in case of
		// transient errors.
		readPersonByName := `
			MATCH (p:Person)
			WHERE p.name = $person_name
			RETURN p.name AS name`
		result, err := tx.Run(readPersonByName, map[string]interface{}{
			"person_name": "Alice",
		})
		if err != nil {
			return nil, err
		}
		// Iterate over the result within the transaction instead of using
		// Collect (just to show how it looks...). Result.Next returns true
		// while a record could be retrieved, in case of error result.Err()
		// will return the error.
		for result.Next() {
			fmt.Printf("Person name: '%s' \n", result.Record().Values[0].(string))
		}
		// Again, return any error back to driver to indicate rollback and
		// retry in case of transient error.
		return nil, result.Err()
	})
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "pong")
		if err != nil {
			panic(err)
		}
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
