package mongodb

func (col dataset) AddLike(owner, rid string) error {
	return updateResourceLike(col.collectionName, owner, rid, 1)
}

func (col dataset) RemoveLike(owner, rid string) error {
	return updateResourceLike(col.collectionName, owner, rid, -1)
}
