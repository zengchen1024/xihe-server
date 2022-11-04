package mongodb

import "sort"

type firstLetterSortData struct {
	level    int
	index    int
	letter   byte
	updateAt int64
}

func firstLetterSortAndPaginate(
	datas []firstLetterSortData, countPerPage, pageNum int,
) []firstLetterSortData {
	i, j, ok := paginate(countPerPage, pageNum, len(datas))
	if !ok {
		return nil
	}

	sort.Slice(datas, func(i, j int) bool {
		a, b := &datas[i], &datas[j]

		if a.level != b.level {
			return a.level > b.level
		}

		if a.letter != b.letter {
			return a.letter < b.letter
		}

		return a.updateAt >= b.updateAt
	})

	return datas[i:j]
}

type updateAtSortData struct {
	id       string
	level    int
	index    int
	updateAt int64
}

func updateAtSortAndPaginate(
	datas []updateAtSortData, countPerPage, pageNum int,
) []updateAtSortData {
	i, j, ok := paginate(countPerPage, pageNum, len(datas))
	if !ok {
		return nil
	}

	sort.Slice(datas, func(i, j int) bool {
		a, b := &datas[i], &datas[j]

		if a.level != b.level {
			return a.level > b.level
		}

		if a.updateAt != b.updateAt {
			return a.updateAt > b.updateAt
		}

		return a.id < b.id
	})

	return datas[i:j]
}

type downloadSortData struct {
	level    int
	index    int
	download int
	updateAt int64
}

func downloadSortAndPaginate(
	datas []downloadSortData, countPerPage, pageNum int,
) []downloadSortData {
	i, j, ok := paginate(countPerPage, pageNum, len(datas))
	if !ok {
		return nil
	}

	sort.Slice(datas, func(i, j int) bool {
		a, b := &datas[i], &datas[j]

		if a.level != b.level {
			return a.level > b.level
		}

		if a.download != b.download {
			return a.download > b.download
		}

		return a.updateAt >= b.updateAt
	})

	return datas[i:j]
}

func paginate(countPerPage, pageNum, total int) (i, j int, ok bool) {
	if total <= 0 {
		return
	}

	if countPerPage <= 0 {
		j = total
		ok = true

		return
	}

	if pageNum > 1 {
		skip := countPerPage * (pageNum - 1)
		if skip >= total {
			return
		}

		i = skip
		total -= skip
	}

	if total >= countPerPage {
		j = i + countPerPage
	} else {
		j = i + total
	}
	ok = true

	return
}
