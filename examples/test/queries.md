

pull all nodes and relationships
```graph
MATCH (n)-[r]->(m) RETURN n, r, m
```

delete all nodes and relationships
```graph
MATCH (n)
DETACH DELETE n
```

