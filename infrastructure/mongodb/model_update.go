package mongodb

func (col model) AddLike(owner, rid string) error {
	return updateResourceLike(col.collectionName, owner, rid, 1)
}

func (col model) RemoveLike(owner, rid string) error {
	return updateResourceLike(col.collectionName, owner, rid, -1)
}
