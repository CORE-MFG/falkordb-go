package main

import (
	"fmt"
	"log"

	"github.com/FalkorDB/falkordb-go"
)

func main() {
	// Connect to FalkorDB
	db, err := falkordb.FromURL("falkor://0.0.0.0:6379")
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer db.Conn.Close()

	// Select a graph
	graph := db.SelectGraph("social")

	// Delete all existing nodes and relationships
	deleteQuery := `
		MATCH (n)
		DETACH DELETE n
	`
	_, err = graph.Query(deleteQuery, nil, nil)
	if err != nil {
		log.Fatal("Failed to delete existing data:", err)
	}

	// Create users, roles and role hierarchy
	q := `
		CREATE 
		// Create base user node
		(userNode:UserType {name:'User'}),

		// Create users
		(john:User {name:'John Doe'}),
		(jane:User {name:'Jane Doe'}),
		(dwight:User {name:'Dwight Schrute'}),
		
		// Create base role node
		(roleNode:RoleType {name:'Role'}),
		
		// Create roles
		(admin:Role {name:'Administrator'}),
		(dev:Role {name:'Developer'}),
		(qa:Role {name:'QA Engineer'}),
		(pm:Role {name:'Project Manager'}),
		(arch:Role {name:'Solutions Architect'}),
		
		// Connect roles to base role node
		(admin)-[:IS_ROLE]->(roleNode),
		(dev)-[:IS_ROLE]->(roleNode),
		(qa)-[:IS_ROLE]->(roleNode),
		(pm)-[:IS_ROLE]->(roleNode),
		(arch)-[:IS_ROLE]->(roleNode),

		// Connect users to base user node
		(john)-[:IS_USER]->(userNode),
		(jane)-[:IS_USER]->(userNode),
		(dwight)-[:IS_USER]->(userNode),
		
		// Create relationships between users and roles
		(john)-[:HAS_ROLE]->(dev),
		(jane)-[:HAS_ROLE]->(qa),
		(dwight)-[:HAS_ROLE]->(admin),
		(dwight)-[:HAS_ROLE]->(arch)
		RETURN john, jane, dwight
	`
	res, err := graph.Query(q, nil, nil)
	if err != nil {
		log.Fatal("Failed to execute query:", err)
	}

	// Print created users
	if res.Next() {
		record := res.Record()
		for i := 0; i < 3; i++ {
			user := record.GetByIndex(i).(*falkordb.Node)
			fmt.Printf("Created user: %s\n", user.GetProperty("name"))
		}
	}

	// Query to verify users and their roles
	q = `
		MATCH (u:User)-[r:HAS_ROLE]->(role:Role)
		RETURN u.name as User, collect(role.name) as Roles
	`
	res, err = graph.Query(q, nil, nil)
	if err != nil {
		log.Fatal("Failed to execute query:", err)
	}

	fmt.Println("\nUsers and their roles:")
	res.PrettyPrint()
}
