package base

type Bucket struct {
	Start int
	End   int
}

type Range struct {
	buckets []*Bucket
}

func NewRange(size int, bucketSize int) *Range {
	r := &Range{}
	r.buckets = make([]*Bucket, 0, size/bucketSize+1)
	if size == 0 {
		return r
	}

	start := 0
	for i := 0; i < size; {
		i += bucketSize
		if i > size {
			i = size
		}
		r.buckets = append(r.buckets, &Bucket{
			Start: start,
			End:   i,
		})
		start = i
	}
	return r
}

func (r *Range) ForEach(f func(b *Bucket) error) error {
	for _, b := range r.buckets {
		err := f(b)
		if err != nil {
			return err
		}
	}
	return nil
}
