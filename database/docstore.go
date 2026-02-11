package database

import (
	gc_docstore "gocloud.dev/docstore"
)

func newDocstoreQuery(docstore_col *gc_docstore.Collection, database_q *Query) *gc_docstore.Query {

	docstore_q := docstore_col.Query()
	return populateDocstoreQuery(docstore_q, database_q)
}

func populateDocstoreQuery(docstore_q *gc_docstore.Query, database_q *Query) *gc_docstore.Query {

	if database_q != nil {

		if database_q.Where != nil {
			
			for _, c := range database_q.Where.Conditions {
				docstore_q = docstore_q.Where(gc_docstore.FieldPath(c.Field), c.Operator, c.Value)
			}
		}
		
		if database_q.OrderBy != nil {
			docstore_q = docstore_q.OrderBy(database_q.OrderBy.Field, database_q.OrderBy.Direction)
		}
		
		if database_q.Offset != nil {
			docstore_q = docstore_q.Offset(*database_q.Offset)
		}
		
		if database_q.Limit != nil {
			docstore_q = docstore_q.Limit(*database_q.Limit)
		}
	}
	
	return docstore_q
}
