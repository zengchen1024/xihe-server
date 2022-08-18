package mongodb

func (col project) AddLike(owner, pid string) error {
	return updateResourceLike(col.collectionName, owner, pid, 1)
}

func (col project) RemoveLike(owner, pid string) error {
	return updateResourceLike(col.collectionName, owner, pid, -1)
}
