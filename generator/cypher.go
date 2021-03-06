package generator

import (
	"fmt"
	"github.com/restra-social/jypher/models"
	"regexp"
	"strings"
)

// CypherGenerator : Cypher Query Generator
type CypherGenerator struct{}

// Generate : This method takes a graph model and generates cypher query
func (c *CypherGenerator) Generate(id string, models map[string]models.Graph, serial []string) (cypher string) {

	// loop through the serial
	for _, term := range serial {

		// search for key in the model in ascending order to generate the query
		if k, ok := models[term]; ok {

			level := k.Nodes.Lebel
			nodeRelName := regexp.MustCompile(`[A-Za-z]+`).FindString(level)

			// Filter Special Label like type which is similar to many resource but has different meaning
			// like Organization Type is different than Claim Type

			if strings.HasPrefix(level, "type") {
				// append source
				level = fmt.Sprintf("%s%s", k.Edges.Source, level)
			}

			pl := len(k.Nodes.Properties)

			node := regexp.MustCompile(`[A-Za-z]+`).FindString(strings.Title(level))
			source := regexp.MustCompile(`[A-Za-z]+`).FindString(k.Edges.Source)

			relation := fmt.Sprintf("%s_%s", strings.ToUpper(source), strings.ToUpper(nodeRelName))

			if k.Nodes.ID != "" {
				cypher += fmt.Sprintf("MERGE (%s:%s {id:'%s'}) SET ", level, node, k.Nodes.ID)

				for i, property := range k.Nodes.Properties {
					for key, val := range property {
						// Using ' ' in value assignment so filter for text contains ''
						filteredVal := strings.Replace(val.(string), "'", "", -1)
						cypher += fmt.Sprintf("%s.%s = '%s'", k.Nodes.Lebel, key, filteredVal)
					}
					if pl > 1 {
						if i < pl-1 {
							cypher += ", "
						}
					}

				}
				if k.Edges.Source != k.Edges.Target { // avoids self loop
					cypher += fmt.Sprintf("MERGE (%s)-[:%s]->(%s)", k.Edges.Source, relation, level)
				}

				cypher += " "

			} else {

				node := regexp.MustCompile(`[A-Za-z]+`).FindString(strings.Title(level))

				len := len(k.Nodes.Properties)

				// If property found then take them for full merge
				if len > 0 {

					cypher += fmt.Sprintf("MERGE (%s:%s { ", level, node)

					for i, property := range k.Nodes.Properties {
						for key, val := range property {
							// Using ' ' in value assignment so filter for text contains ''
							filteredVal := strings.Replace(val.(string), "'", "", -1)
							cypher += fmt.Sprintf("%s:'%s'", key, filteredVal)

							// skip comma for last property
							if i != len-1 {
								cypher += fmt.Sprint(",")
							}
						}

						if i == len-1 {
							cypher += fmt.Sprint(" }) ")
						}
					}

					// append the id to each node
					cypher += fmt.Sprintf("SET %s._id = '%s' ", level, id)

					// Add Relation
					cypher += fmt.Sprintf("MERGE (%s)-[:%s]->(%s)", k.Edges.Source, relation, level)

					// Add Gap
					cypher += " "

				} else {

					// for those nodes who doesn't have any properties
					cypher += fmt.Sprintf("MERGE (%s:%s)", level, node)

					/*cypher += fmt.Sprintf("MERGE (%s:%s) SET ", level, node)

					for _, property := range k.Nodes.Properties {
						for key, val := range property {
							cypher += fmt.Sprintf("%s.%s = '%s', ", k.Nodes.Lebel, key, val)
						}
					}

					// append the id to each node
					cypher += fmt.Sprintf("%s._id = '%s' ", k.Nodes.Lebel, id)*/

					cypher += fmt.Sprintf("MERGE (%s)-[:%s]->(%s)", k.Edges.Source, relation, k.Edges.Target)

					cypher += " "
				}
			}

		}
	}

	return cypher
}

// #todo #fix
// fix MERGE (patient:Patient {id:'34876259-35cd-497c-a932-94baaaeb555c'}) SET patient.reference = 'urn:uuid:34876259-35cd-497c-a932-94baaaeb555c'
// fix
