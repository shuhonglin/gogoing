package util

type ring struct {
	count, i int
	data     []Resource
}

func (r *ring) Size() int {
	return r.count
}

func (r *ring) Empty() bool {
	return r.count == 0
}

func (r *ring) Peek() Resource {
	return r.data[r.i]
}

func (r *ring) Enqueue(res Resource) {
	if r.count>=len(r.data) {
		r.grow(2*r.count+1)
	}
	r.data[(r.i+r.count)%len(r.data)] = res
	r.count++
}

func (r *ring) Dequeue()(res Resource) {
	res = r.Peek()
	r.count,r.i = r.count-1,(r.i+1)%len(r.data)
	return
}

func (r *ring) grow(newSize int) {
	newData := make([]Resource, newSize)
	n := copy(newData, r.data[r.i:])
	copy(newData[n:], r.data[:r.count-n])

	r.i=0
	r.data = newData
}
